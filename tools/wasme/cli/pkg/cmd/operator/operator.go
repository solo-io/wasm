package operator

import (
	"context"
	"time"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/skv2/pkg/ezkube"
	cachedeployment "github.com/solo-io/wasm/tools/wasme/cli/pkg/cache"
	"github.com/solo-io/wasm/tools/wasme/cli/pkg/deploy/istio"
	"github.com/solo-io/wasm/tools/wasme/cli/pkg/operator"
	v1 "github.com/solo-io/wasm/tools/wasme/cli/pkg/operator/api/wasme.io/v1"
	"github.com/solo-io/wasm/tools/wasme/cli/pkg/operator/api/wasme.io/v1/controller"
	"github.com/solo-io/wasm/tools/wasme/pkg/version"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
	"sigs.k8s.io/controller-runtime/pkg/log"
	zaputil "sigs.k8s.io/controller-runtime/pkg/log/zap"

	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// set this log level via viper flag
type flagSetLogLevel struct {
	zapcore.Level
}

func (f flagSetLogLevel) Type() string {
	return "zap log level"
}

type operatorOpts struct {
	cache        istio.Cache
	logLevel     flagSetLogLevel
	cacheTimeout time.Duration
}

func OperatorCmd(ctx *context.Context) *cobra.Command {
	var opts operatorOpts

	cmd := &cobra.Command{
		Use:   "operator [--cache-name=<cache name>] [--cache-namespace=<cache namespace>]",
		Short: "Run the Wasme Kubernetes Operator",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runOperator(*ctx, opts)
		},
		Hidden: true,
	}

	cmd.Flags().StringVar(&opts.cache.Name, "cache-name", cachedeployment.CacheName, "name of resources for the wasm image cache server")
	cmd.Flags().StringVar(&opts.cache.Namespace, "cache-namespace", cachedeployment.CacheNamespace, "namespace of resources for the wasm image cache server")
	cmd.Flags().Var(&opts.logLevel, "log-level", "the logging level to use")
	cmd.Flags().DurationVar(&opts.cacheTimeout, "cache-timeout", time.Minute, "the length of time to wait for the server-side filter cache to pull the filter image before giving up with an error. set to 0 to skip the check entirely (note, this may produce a known race condition).")

	return cmd
}

func runOperator(ctx context.Context, opts operatorOpts) error {
	zapLevel := zap.NewAtomicLevel()
	zapLevel.SetLevel(opts.logLevel.Level)
	log.SetLogger(zaputil.New(
		zaputil.Level(&zapLevel),
	))

	contextutils.LoggerFrom(ctx).Infof("started wasme version %v", version.Version)
	// get local kubeconfig
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	// create manager
	mgr, err := manager.New(cfg, manager.Options{
		Namespace:          "", // watch all namespaces
		MetricsBindAddress: ":9091",
	})
	if err != nil {
		return err
	}
	// add CRDs to scheme
	if err := v1.AddToScheme(mgr.GetScheme()); err != nil {
		return err
	}
	// create controller
	ctl := controller.NewFilterDeploymentEventWatcher("wasme", mgr)

	// kube client
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}

	// ezkube client wrapper
	client := ezkube.NewEnsurer(ezkube.NewRestClient(mgr))

	// create handler
	handler := operator.NewFilterDeploymentHandler(
		ctx,
		kubeClient,
		client,
		opts.cache,
		opts.cacheTimeout,
	)

	eg := &errgroup.Group{}
	eg.Go(func() error {
		return ctl.AddEventHandler(ctx, handler)
	})
	eg.Go(func() error {
		return mgr.Start(ctx.Done())
	})
	return eg.Wait()
}
