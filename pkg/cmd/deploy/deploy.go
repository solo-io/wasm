package deploy

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var log = logrus.StandardLogger()

func DeployCmd() *cobra.Command {
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
			if len(args) == 0 {
				return fmt.Errorf("invalid number of arguments")
			}
			opts.filter.Image = args[0]
			if opts.filter.ID == "" {
				return errors.Errorf("--id cannot be empty")
			}
			return nil
		},
	}
	opts.addToFlags(cmd.PersistentFlags())

	cmd.AddCommand(
		deployGlooCmd(opts),
		deployLocalCmd(opts),
	)

	return cmd
}

func makeDeployCommand(opts *options, provider, use, short, long string, minArgs int, addFlags ...func(flags *pflag.FlagSet)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(minArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.providerType = provider
			return runDeploy(opts)
		},
	}

	for _, f := range addFlags {
		f(cmd.PersistentFlags())
	}

	return cmd
}

func deployGlooCmd(opts *options) *cobra.Command {
	use := "gloo <image> --id=<unique name> [--config=<inline string>] [--root-id=<root id>] [--namespaces <comma separated namespaces>] [--labels <key1=val1,key2=val2>]"
	short := "Deploy an Envoy WASM Filter to the Gloo Gateway Proxies (Envoy)."
	long := `Deploys an Envoy WASM Filter to Gloo Gateway Proxies.

wasme uses the Gloo Gateway CR to pull and run wasm filters.

Use --namespaces to constrain the namespaces of Gateway CRs to update.

Use --labels to use a match Gateway CRs by label.
`
	return makeDeployCommand(opts,
		Provider_Gloo,
		use,
		short,
		long,
		1,
		opts.glooOpts.addToFlags,
	)
}

func deployLocalCmd(opts *options) *cobra.Command {
	use := "envoy <image> --id=<unique id> [--config=<inline string>] [--root-id=<root id>] --in=<input config file> --out=<output config file> --filter <path to filter wasm> [--use-json]"
	short := "Configure a local instance of Envoy to run a WASM Filter."
	long := `
Unlike ` + "`" + `wasme deploy gloo` + "`" + ` and ` + "`" + `wasme deploy istio` + "`" + `, ` + "`" + `wasme deploy envoy` + "`" + ` only outputs the Envoy configuration required to run the filter with Envoy.

Launch Envoy using the output configuration to run the wasm filter.
`
	return makeDeployCommand(opts,
		Provider_Envoy,
		use,
		short,
		long,
		1,
		opts.localOpts.addToFlags,
	)
}

func runDeploy(opts *options) error {
	deployer, err := makeDeployer(opts)
	if err != nil {
		return err
	}

	if opts.remove {
		return deployer.RemoveFilter(&opts.filter)
	}

	return deployer.ApplyFilter(&opts.filter)
}
