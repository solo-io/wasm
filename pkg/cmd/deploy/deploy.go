package deploy

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var log = logrus.StandardLogger()

type deployOptions struct {
	provider        string
	operation       string
	dryRun          bool
	injectName      string
	injectNamespace string
}

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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("invalid number of arguments")
			}
			opts.filter.Image = args[0]
			if opts.filter.ID == "" {
				return errors.Errorf("--id cannot be empty")
			}
			return runDeploy(opts)
		},
	}
	opts.addToFlags(cmd.PersistentFlags())

	cmd.AddCommand(
		deployGlooCmd(opts),
	)

	return cmd
}

func deployGlooCmd(opts *options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gloo <image> --id=<unique name> [--config=<inline string>] [--root-id=<root id>] [--namespaces <comma separated namespaces>] [--labels <key1=val1,key2=val2>",
		Short: "Deploy an Envoy WASM Filter to the Gloo Gateway Proxies (Envoy).",
		Long: `Deploys an Envoy WASM Filter to Gloo Gateway Proxies.

wasme uses the Gloo Gateway CR to pull and run wasm filters.

Use --namespaces to constrain the namespaces of Gateway CRs to update.

Use --labels to use a match Gateway CRs by label.
`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.providerType = Provider_Gloo
			return runDeploy(opts)
		},
	}

	opts.glooOpts.addToFlags(cmd.Flags())

	return cmd
}

func deployLocalCmd(opts *options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "envoy <>",
		Short: "Configure a local instance of Envoy to run a WASM Filter.",
		Long: `
Unlike ` + "`" + `wasme deploy gloo` + "`" + ` and ` + "`" + `wasme deploy istio` + "`" + `, ` + "`" + `wasme deploy envoy` + "`" + ` only outputs the Envoy configuration required to run the filter with Envoy.

Launch Envoy using the output configuration to run the wasm filter.

`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.providerType = Provider_Envoy
			return runDeploy(opts)
		},
	}

	opts.localOpts.addToFlags(cmd.Flags())

	return cmd
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
