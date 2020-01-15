package main_test

import (
	"io/ioutil"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/autopilot/cli/pkg/utils"
	"github.com/solo-io/autopilot/codegen/util"
	"github.com/solo-io/autopilot/test"
	"github.com/solo-io/go-utils/kubeutils"
	"github.com/solo-io/go-utils/randutils"
	kubehelp "github.com/solo-io/go-utils/testutils/kube"
	"github.com/solo-io/wasme/pkg/operator/api/wasme.io/v1/clientset/versioned"
	"github.com/solo-io/wasme/pkg/operator/api/wasme.io/v1/controller"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	zaputil "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func applyFile(file string) error {
	path := filepath.Join(util.MustGetThisDir(), "manifests", file)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return utils.KubectlApply(b)
}

func deleteFile(file string) error {
	path := filepath.Join(util.MustGetThisDir(), "manifests", file)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return utils.KubectlDelete(b)
}

var _ = Describe("AutopilotGenerate", func() {
	var (
		ns        string
		kube      kubernetes.Interface
		clientset versioned.Interface
		logLevel  = zap.NewAtomicLevel()
	)
	BeforeEach(func() {
		logLevel.SetLevel(zap.DebugLevel)
		log.SetLogger(zaputil.New(
			zaputil.Level(&logLevel),
		))
		log.Log.Info("test")
		err := applyFile("things.test.io_v1_crds.yaml")
		Expect(err).NotTo(HaveOccurred())
		ns = randutils.RandString(4)
		kube = kubehelp.MustKubeClient()
		err = kubeutils.CreateNamespacesInParallel(kube, ns)
		Expect(err).NotTo(HaveOccurred())
		clientset, err = versioned.NewForConfig(test.MustConfig())
		Expect(err).NotTo(HaveOccurred())
	})
	AfterEach(func() {
		err := deleteFile("things.test.io_v1_crds.yaml")
		Expect(err).NotTo(HaveOccurred())
		err = kubeutils.DeleteNamespacesInParallelBlocking(kube, ns)
		Expect(err).NotTo(HaveOccurred())
	})

	It("runs the wasme operator", func() {
		cfg, err := config.GetConfig()
		Expect(err).NotTo(HaveOccurred())

		mgr, err := manager.New(cfg, manager.Options{
			Namespace:        "",
			EventBroadcaster: record.NewBroadcaster(),
		})
		Expect(err).NotTo(HaveOccurred())

		ctl, err := controller.NewFilterDeploymentController("wasme", mgr)
		Expect(err).NotTo(HaveOccurred())

		err = ctl.AddEventHandler(&controller.FilterDeploymentEventHandlerFuncs{
			OnCreate: nil,
			OnUpdate: nil,
			OnDelete: nil,
		})
		Expect(err).NotTo(HaveOccurred())
	})
})
