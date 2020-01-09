package deploy

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func UndeployCmd(ctx *context.Context) *cobra.Command {
	opts := &options{
		remove: true,
	}
	cmd := &cobra.Command{
		Use:   "undeploy gloo|istio|envoy --id=<unique id>",
		Short: "Remove a deployed Envoy WASM Filter from the data plane (Envoy proxies).",
		Long: `Removes a deployed Envoy WASM Filter from Envoy instances.

`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.filter.ID == "" {
				return errors.Errorf("--id cannot be empty")
			}
			return nil
		},
	}
	opts.addDryRunToFlags(cmd.PersistentFlags())
	opts.addIdToFlags(cmd.PersistentFlags())

	cmd.AddCommand(
		undeployGlooCmd(ctx, opts),
		undeployIstioCmd(ctx, opts),
		undeployLocalCmd(ctx, opts),
	)

	return cmd
}

func undeployGlooCmd(ctx *context.Context, opts *options) *cobra.Command {
	use := "gloo --id=<unique name>"
	short := "Remove an Envoy WASM Filter from the Gloo Gateway Proxies (Envoy)."
	long := `wasme uses the Gloo Gateway CR to pull and run wasm filters.

Use --namespaces to constrain the namespaces of Gateway CRs to update.

Use --labels to use a match Gateway CRs by label.
`
	return makeDeployCommand(ctx, opts,
		Provider_Gloo,
		use,
		short,
		long,
		0,
		opts.glooOpts.addToFlags,
	)
}

func undeployIstioCmd(ctx *context.Context, opts *options) *cobra.Command {
	use := "istio --id=<unique name> --namespace=<deployment namespace> [--name=<deployment name>]"
	short := "Remove an Envoy WASM Filter from the Istio Sidecar Proxies (Envoy)."
	long := `wasme uses the Istio EnvoyFilter CR to pull and run wasm filters.

Use --namespace to target workload(s) in a the namespaces of Gateway CRs to update.
Use --name to target a specific workload (deployment or daemonset) in the target namespace. If unspecified, all deployments 
in the namespace will be targeted.
`
	return makeDeployCommand(ctx, opts,
		Provider_Istio,
		use,
		short,
		long,
		0,
		opts.providerOptions.istioOpts.addToFlags,
	)
}

func undeployLocalCmd(ctx *context.Context, opts *options) *cobra.Command {
	use := "envoy --id=<unique name>"
	short := "Remove an Envoy WASM Filter from the Envoy listeners."
	long := `wasme removes the deployed filter matching the given id. 
`
	return makeDeployCommand(ctx, opts,
		Provider_Envoy,
		use,
		short,
		long,
		0,
		opts.localOpts.addFilesToFlags,
	)
}
