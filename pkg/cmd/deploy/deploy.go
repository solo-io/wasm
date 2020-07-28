package deploy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gogo/protobuf/types"
	"github.com/solo-io/wasme/pkg/deploy/local"
	v1 "github.com/solo-io/wasme/pkg/operator/api/wasme.io/v1"
	"github.com/solo-io/wasme/pkg/store"

	corev1 "k8s.io/api/core/v1"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	cachedeployment "github.com/solo-io/wasme/pkg/cache"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var log = logrus.StandardLogger()

func DeployCmd(ctx *context.Context, parentPreRun func(cmd *cobra.Command, args []string)) *cobra.Command {
	opts := &options{}
	cmd := &cobra.Command{
		Use:   "deploy gloo|istio|envoy <image> --id=<unique id> [--config=<inline string>] [--root-id=<root id>]",
		Short: "Deploy an Envoy WASM Filter to the data plane (Envoy proxies).",
		Long: `Deploys an Envoy WASM Filter to Envoy instances.

You must provide a value for --id which will become the unique ID of the deployed filter. When using --provider=istio, the ID must be a valid Kubernetes resource name.

You must specify --root-id unless a default root id is provided in the image configuration. Use --root-id to select the filter to run if the wasm image contains more than one filter.

`,
		Args: cobra.MinimumNArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			parentPreRun(cmd, args)
			if len(args) == 0 {
				return fmt.Errorf("invalid number of arguments")
			}
			opts.filter.Image = args[0]
			return nil
		},
	}
	opts.addToFlags(cmd.PersistentFlags())

	cmd.AddCommand(
		deployGlooCmd(ctx, opts),
		deployIstioCmd(ctx, opts),
		deployLocalCmd(ctx, opts),
	)

	return cmd
}

func makeDeployCommand(ctx *context.Context, opts *options, provider, use, short, long string, minArgs int, addFlags ...func(flags *pflag.FlagSet)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(minArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.filter.Id == "" {
				return errors.Errorf("--id cannot be empty")
			}
			opts.providerType = provider
			// If we were passed a config via CLI flag, default config type to StringValue
			if opts.filterConfig != "" {
				opts.filter.Config = &types.Any{
					TypeUrl: "type.googleapis.com/google.protobuf.StringValue",
					Value:   []byte(opts.filterConfig),
				}
			}
			return runDeploy(*ctx, opts)
		},
	}

	opts.addToFlags(cmd.PersistentFlags())

	for _, f := range addFlags {
		f(cmd.PersistentFlags())
	}

	return cmd
}

func deployGlooCmd(ctx *context.Context, opts *options) *cobra.Command {
	use := "gloo <image> --id=<unique name> [--config=<inline string>] [--root-id=<root id>] [--namespaces <comma separated namespaces>] [--labels <key1=val1,key2=val2>]"
	short := "Deploy an Envoy WASM Filter to the Gloo Gateway Proxies (Envoy)."
	long := `Deploys an Envoy WASM Filter to Gloo Gateway Proxies.

wasme uses the Gloo Gateway CR to pull and run wasm filters.

Use --namespaces to constrain the namespaces of Gateway CRs to update.

Use --labels to use a match Gateway CRs by label.
`
	return makeDeployCommand(ctx, opts,
		Provider_Gloo,
		use,
		short,
		long,
		1,
		opts.glooOpts.addToFlags,
	)
}

func deployIstioCmd(ctx *context.Context, opts *options) *cobra.Command {
	use := "istio <image> --id=<unique name> [--config=<inline string>] [--root-id=<root id>] [--namespaces <comma separated namespaces>] [--name deployment-name]"
	short := "Deploy an Envoy WASM Filter to Istio Sidecar Proxies (Envoy)."
	long := `Deploy an Envoy WASM Filter to Istio Sidecar Proxies (Envoy).

wasme uses the EnvoyFilter Istio Custom Resource to pull and run wasm filters.
wasme deploys a server-side cache component which runs in cluster and pulls filter images.

If --name is not provided, all deployments in the targeted namespace will attach the filter.

Note: currently only Istio 1.5.x is supported.
`
	cmd := makeDeployCommand(ctx, opts,
		Provider_Istio,
		use,
		short,
		long,
		1,
		opts.istioOpts.addToFlags,
		opts.cacheOpts.addToFlags,
	)

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		cacheDeployer := cachedeployment.NewDeployer(
			helpers.MustKubeClient(),
			opts.cacheOpts.namespace,
			opts.cacheOpts.name,
			opts.cacheOpts.imageRepo,
			opts.cacheOpts.imageTag,
			opts.cacheOpts.customArgs,
			corev1.PullPolicy(opts.cacheOpts.pullPolicy),
		)

		return cacheDeployer.EnsureCache()
	}

	return cmd
}

func deployLocalCmd(ctx *context.Context, opts *options) *cobra.Command {
	use := "envoy <image> [--config=<filter config>] [--bootstrap=<custom envoy bootstrap file>] [--envoy-image=<custom envoy image>]"
	short := "Run Envoy locally in Docker and attach a WASM Filter."
	long := `
This command runs Envoy locally in docker using a static bootstrap configuration which includes 
the specified WASM filter image. 

The bootstrap can be generated from an internal default or a modified config provided by the user with --bootstrap.

The generated bootstrap config can be output to a file with --out. If using this option, Envoy will not be started locally.
`

	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.filter.Image = args[0]
			return runLocalEnvoy(*ctx, opts.filter, opts.localOpts)
		},
	}

	opts.localOpts.addToFlags(cmd.Flags())

	return cmd
}

func runDeploy(ctx context.Context, opts *options) error {
	deployer, err := makeDeployer(ctx, opts)
	if err != nil {
		return err
	}

	if opts.remove {
		return deployer.RemoveFilter(&opts.filter)
	}

	return deployer.ApplyFilter(&opts.filter)
}

func runLocalEnvoy(ctx context.Context, filter v1.FilterSpec, opts localOpts) error {
	in, err := func() (io.ReadCloser, error) {
		switch opts.infile {
		case "-":
			// use stdin
			return os.Stdin, nil
		case "":
			// use default config
			return ioutil.NopCloser(bytes.NewBuffer([]byte(local.BasicEnvoyConfig))), nil
		default:
			// read file
			return os.Open(opts.infile)
		}
	}()
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := func() (io.WriteCloser, error) {
		switch opts.outfile {
		case "-":
			// use stdout
			return os.Stdout, nil
		case "":
			// use default config
			return nil, nil
		default:
			// write file
			return os.Create(opts.outfile)
		}
	}()
	if err != nil {
		return err
	}
	if out != nil {
		defer out.Close()
	}

	parseArgs := func(argStr string) []string {
		if argStr == "" {
			return nil
		}
		return strings.Split(argStr, " ")
	}

	runner := &local.Runner{
		Ctx:              ctx,
		Input:            in,
		Output:           out,
		Store:            store.NewStore(opts.storageDir),
		DockerRunArgs:    parseArgs(opts.dockerRunArgs),
		EnvoyArgs:        parseArgs(opts.envoyArgs),
		EnvoyDockerImage: opts.envoyDockerImage,
	}

	return runner.RunFilter(&filter)
}
