package operator

import (
	"context"

	"github.com/solo-io/autopilot/pkg/ezkube"
	"github.com/solo-io/go-utils/contextutils"
	cachedeployment "github.com/solo-io/wasme/pkg/cache"
	"github.com/solo-io/wasme/pkg/deploy/istio"
	"github.com/solo-io/wasme/pkg/operator"
	"github.com/solo-io/wasme/pkg/operator/api/wasme.io/v1/controller"
	"github.com/solo-io/wasme/pkg/version"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	cache    istio.Cache
	logLevel flagSetLogLevel
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

	// create controller
	ctl, err := controller.NewFilterDeploymentController("wasme", mgr)
	if err != nil {
		return err
	}

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
	)

	// register the handler to the controller
	if err := ctl.AddEventHandler(ctx, handler); err != nil {
		return err
	}

	// start the manager
	return mgr.Start(ctx.Done())
}
