package deploy

import (
	"context"
	"github.com/pkg/errors"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	"github.com/solo-io/wasme/pkg/cmd/opts"
	"github.com/solo-io/wasme/pkg/deploy"
	"github.com/solo-io/wasme/pkg/deploy/gloo"
	"github.com/solo-io/wasme/pkg/deploy/local"
	"github.com/solo-io/wasme/pkg/pull"
	"github.com/solo-io/wasme/pkg/resolver"
	"github.com/spf13/pflag"
	"io"
	"os"
)

type options struct {
	// filter to deploy
	filter deploy.Filter

	// deployment implementation
	providerOptions

	// login
	opts.GeneralOptions

	// print yaml only
	dryRun bool

	// remove a deployed filter instead of deploying
	remove bool
}


func (opts *options) addToFlags(flags *pflag.FlagSet) {

	flags.StringVarP(&opts.filter.Config, "config", "", "", "optional config that will be passed to the filter. accepts an inline string.")
	flags.StringVarP(&opts.filter.RootID, "root-id", "", "", "optional root ID used to bind the filter at the Envoy level. this value is normally read from the filter image directly.")
	flags.StringVar(&opts.filter.ID, "id", "", "path to the output .wasm file. Nonexistent directories will be created.")
	flags.BoolVarP(&opts.dryRun, "dry-run", "", false, "print output any configuration changes to stdout rather than applying them to the target file / kubernetes cluster")
}


type providerOptions struct {
	providerType string

	glooOpts  glooOpts
	localOpts localOpts
}

type glooOpts struct {
	selector gloo.Selector
}


func (opts *glooOpts) addToFlags(flags *pflag.FlagSet) {
	flags.StringSliceVarP(&opts.selector.Namespaces, "namespaces", "n", nil, "deploy the filter to selected Gateway resource in the given namespaces. if none provided, Gateways in all namespaces will be selected.")
	flags.StringToStringVarP(&opts.selector.GatewayLabels, "labels", "l", nil, "select deploy the filter to selected Gateway resource in the given namespaces. if none provided, Gateways in all namespaces will be selected.")
}


type localOpts struct {
	infile        string
	outfile       string
	filterPath    string
	useJsonConfig bool
}

func (opts *localOpts) addToFlags(flags *pflag.FlagSet) {
	flags.StringSliceVarP(&opts.selector.Namespaces, "namespaces", "n", nil, "deploy the filter to selected Gateway resource in the given namespaces. if none provided, Gateways in all namespaces will be selected.")
	flags.StringToStringVarP(&opts.selector.GatewayLabels, "labels", "l", nil, "select deploy the filter to selected Gateway resource in the given namespaces. if none provided, Gateways in all namespaces will be selected.")
}

const (
	Provider_Gloo  = "gloo"
	Provider_Istio = "istio"
	Provider_Envoy = "envoy"
)

var SupportedProviders = []string{
	Provider_Gloo,
	Provider_Istio,
	Provider_Envoy,
}

func (opts options) makeProvider(ctx context.Context) (deploy.Provider, error) {
	switch opts.providerType {
	case Provider_Gloo:
		var gwClient gatewayv1.GatewayClient
		if opts.dryRun {
			gwClient = newDryRunGatewayClient(os.Stdout)
		} else {
			gwClient = helpers.MustGatewayV2Client()
		}

		return &gloo.Provider{
			Ctx:           ctx,
			GatewayClient: gwClient,
			Selector:      opts.glooOpts.selector,
		}, nil
	case Provider_Istio:
		return nil, errors.Errorf("istio currently not supported")
	case Provider_Envoy:

		var in io.Reader
		if opts.localOpts.infile == "-" {
			// use stdin
			in = os.Stdin
		} else {
			f, err := os.Open(opts.localOpts.infile)
			if err != nil {
				return nil, err
			}
			in = f
		}

		var out io.Writer
		if opts.localOpts.infile == "-" {
			// use stdout
			out = os.Stdout
		} else {
			f, err := os.Open(opts.localOpts.outfile)
			if err != nil {
				return nil, err
			}
			out = f
		}

		return &local.Provider{
			Ctx:           ctx,
			Input:         in,
			Output:        out,
			FilterPath:    opts.localOpts.filterPath,
			UseJsonConfig: opts.localOpts.useJsonConfig,
		}, nil
	}

	return nil, nil
}

func makeDeployer(opts *options) (*deploy.Deployer, error) {
	resolver, _ := resolver.NewResolver(opts.Username, opts.Password, opts.Insecure, opts.PlainHTTP, opts.Configs...)
	puller := pull.NewPuller(resolver)

	ctx := context.Background()

	provider, err := opts.makeProvider(ctx)
	if err != nil {
		return nil, err
	}
	return &deploy.Deployer{
		Ctx:      ctx,
		Puller:   puller,
		Provider: provider,
	}, nil
}
