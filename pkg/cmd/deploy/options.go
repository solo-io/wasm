package deploy

import (
	"context"
	"os"
	"time"

	"github.com/solo-io/wasme/pkg/deploy/local"
	corev1 "k8s.io/api/core/v1"

	"github.com/solo-io/skv2/pkg/ezkube"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/pkg/errors"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	"github.com/solo-io/go-utils/kubeutils"
	cachedeployment "github.com/solo-io/wasme/pkg/cache"
	"github.com/solo-io/wasme/pkg/cmd/opts"
	"github.com/solo-io/wasme/pkg/deploy"
	"github.com/solo-io/wasme/pkg/deploy/gloo"
	"github.com/solo-io/wasme/pkg/deploy/istio"
	v1 "github.com/solo-io/wasme/pkg/operator/api/wasme.io/v1"
	"github.com/solo-io/wasme/pkg/pull"
	"github.com/solo-io/wasme/pkg/resolver"
	"github.com/spf13/pflag"
)

type options struct {
	// filter to deploy
	filter v1.FilterSpec

	// configuration string for filter
	filterConfig string

	// deployment implementation
	providerOptions

	// login
	opts.AuthOptions

	// print yaml only
	dryRun bool

	// remove a deployed filter instead of deploying
	remove bool
}

func (opts *options) addToFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&opts.filterConfig, "config", "", "", "optional config that will be passed to the filter. accepts an inline string.")
	flags.StringVarP(&opts.filter.RootID, "root-id", "", "", "optional root ID used to bind the filter at the Envoy level. this value is normally read from the filter image directly.")
	opts.addIdToFlags(flags)
}

func (opts *options) addDryRunToFlags(flags *pflag.FlagSet) {
	flags.BoolVarP(&opts.dryRun, "dry-run", "", false, "print output any configuration changes to stdout rather than applying them to the target file / kubernetes cluster")
}

func (opts *options) addIdToFlags(flags *pflag.FlagSet) {
	flags.StringVar(&opts.filter.Id, "id", "", "unique id for naming the deployed filter. this is used for logging as well as removing the filter. when running wasme deploy istio, this name must be a valid Kubernetes resource name.")
}

type providerOptions struct {
	providerType string

	glooOpts  glooOpts
	localOpts localOpts
	istioOpts istioOpts

	cacheOpts cacheOpts
}

type glooOpts struct {
	selector gloo.Selector
}

func (opts *glooOpts) addToFlags(flags *pflag.FlagSet) {
	flags.StringSliceVarP(&opts.selector.Namespaces, "namespaces", "n", nil, "deploy the filter to selected Gateway resource in the given namespaces. if none provided, Gateways in all namespaces will be selected.")
	flags.StringToStringVarP(&opts.selector.GatewayLabels, "labels", "l", nil, "select deploy the filter to selected Gateway resource in the given namespaces. if none provided, Gateways in all namespaces will be selected.")
}

type istioOpts struct {
	workload       istio.Workload
	istioNamespace string
	cacheTimeout   time.Duration

	puller pull.ImagePuller // set by load
}

func (opts *istioOpts) addToFlags(flags *pflag.FlagSet) {
	flags.StringToStringVarP(&opts.workload.Labels, "labels", "l", nil, "labels of the deployment or daemonset into which to inject the filter. if not set, will apply to all workloads in the target namespace")
	flags.StringVarP(&opts.workload.Namespace, "namespace", "n", "default", "namespace of the workload(s) to inject the filter.")
	flags.StringVarP(&opts.workload.Kind, "workload-type", "t", istio.WorkloadTypeDeployment, "type of workload into which the filter should be injected. possible values are "+istio.WorkloadTypeDeployment+" or "+istio.WorkloadTypeDaemonSet)
	flags.StringVarP(&opts.istioNamespace, "istio-namespace", "", "istio-system", "the namespace where the Istio control plane is installed")
	flags.DurationVar(&opts.cacheTimeout, "cache-timeout", time.Minute, "the length of time to wait for the server-side filter cache to pull the filter image before giving up with an error. set to 0 to skip the check entirely (note, this may produce a known race condition).")
}

type cacheOpts struct {
	name       string
	namespace  string
	imageRepo  string
	imageTag   string
	customArgs []string
	pullPolicy string
}

func (opts *cacheOpts) addToFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&opts.name, "cache-name", "", cachedeployment.CacheName, "name of resources for the wasm image cache server")
	flags.StringVarP(&opts.namespace, "cache-namespace", "", cachedeployment.CacheNamespace, "namespace of resources for the wasm image cache server")
	flags.StringVarP(&opts.imageRepo, "cache-repo", "", cachedeployment.CacheImageRepository, "name of the image repository to use for the cache server daemonset")
	flags.StringVarP(&opts.imageTag, "cache-tag", "", cachedeployment.CacheImageTag, "image tag to use for the cache server daemonset")
	flags.StringSliceVarP(&opts.customArgs, "cache-custom-command", "", nil, "custom command to provide to the cache server image")
	flags.StringVarP(&opts.pullPolicy, "cache-image-pull-policy", "", string(corev1.PullIfNotPresent), "image pull policy for the cache server daemonset. see https://kubernetes.io/docs/concepts/containers/images/")
}

type localOpts struct {
	infile           string
	outfile          string
	storageDir       string
	dockerRunArgs    string
	envoyArgs        string
	envoyDockerImage string
}

func (opts *localOpts) addToFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&opts.envoyDockerImage, "envoy-image", "e", local.DefaultEnvoyImage, "Name of the Docker image containing the Envoy binary")
	flags.StringVarP(&opts.infile, "bootstrap", "b", "", "Path to an Envoy bootstrap config. If set, `wasme deploy envoy` will run Envoy locally using the provided configuration file. Set -in=- to use stdin. If empty, will use a default configuration template with a single route to `jsonplaceholder.typicode.com`.")
	flags.StringVarP(&opts.outfile, "out", "", "", "If set, write the modified Envoy configuration to this file instead of launching Envoy. Set -out=- to use stdout.")
	flags.StringVar(&opts.storageDir, "store", "", "Set the path to the local storage directory for wasm images. Defaults to $HOME/.wasme/store")
	flags.StringVar(&opts.dockerRunArgs, "docker-run-args", "", "Set to provide additional args to the `docker run` command used to launch Envoy. Ignored if --out is set.")
	flags.StringVar(&opts.envoyArgs, "envoy-run-args", "", "Set to provide additional args to the `envoy` command used to launch Envoy. Ignored if --out is set.")
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

func (opts *options) makeProvider(ctx context.Context) (deploy.Provider, error) {
	switch opts.providerType {
	case Provider_Gloo:
		var gwClient gatewayv1.GatewayClient
		if opts.dryRun {
			gwClient = newDryRunGatewayClient(os.Stdout)
		} else {
			gwClient = helpers.MustGatewayClient()
		}

		return &gloo.Provider{
			Ctx:           ctx,
			GatewayClient: gwClient,
			Selector:      opts.glooOpts.selector,
		}, nil
	case Provider_Istio:
		if opts.dryRun {
			return nil, errors.Errorf("dry-run not currenty supported for istio")
		}

		cfg, err := kubeutils.GetConfig("", "")
		if err != nil {
			return nil, err
		}

		kubeClient, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			return nil, err
		}

		mgr, err := manager.New(cfg, manager.Options{})
		if err != nil {
			return nil, err
		}

		go func() {
			err := mgr.Start(ctx.Done())
			if err != nil {
				log.Fatalf("failed to start kubernetes dynamic client")
			}
		}()

		return istio.NewProvider(
			ctx,
			kubeClient,
			ezkube.NewEnsurer(ezkube.NewRestClient(mgr)),
			opts.istioOpts.puller,
			opts.istioOpts.workload,
			istio.Cache{
				Name:      opts.cacheOpts.name,
				Namespace: opts.cacheOpts.namespace,
			},
			nil, // no parent object when using CLI
			nil, // no callback when using CLI
			opts.istioOpts.istioNamespace,
			opts.istioOpts.cacheTimeout,
		)
	}

	return nil, nil
}

func makeDeployer(ctx context.Context, opts *options) (*deploy.Deployer, error) {
	resolver, _ := resolver.NewResolver(opts.Username, opts.Password, opts.Insecure, opts.PlainHTTP, opts.CredentialsFiles...)
	puller := pull.NewPuller(resolver)

	// set istio puller
	opts.istioOpts.puller = puller

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
