package istio

import (
	"github.com/pkg/errors"
	"github.com/solo-io/wasme/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	pilotDeploymentName   = "istiod"
	defaultIstioNamespace = "istio-system"
	pilotContainerName    = "discovery"
)

type VersionInspector interface {
	GetIstioVersion() (string, error)
}

type versionInspector struct {
	istioNamespace string
	kube           kubernetes.Interface
}

func (i *versionInspector) GetIstioVersion() (string, error) {
	istioNamespace := i.istioNamespace
	if istioNamespace == "" {
		istioNamespace = defaultIstioNamespace
	}
	pilotDeployment, err := i.kube.AppsV1().Deployments(istioNamespace).Get(pilotDeploymentName, metav1.GetOptions{})
	if err != nil {
		return "", nil
	}
	var pilotImage string
	for _, container := range pilotDeployment.Spec.Template.Spec.Containers {
		if container.Name == pilotContainerName {
			pilotImage = container.Image
			break
		}
	}
	if pilotImage == "" {
		return "", errors.Errorf("did not find container named %s on pilot deployment", pilotContainerName)
	}

	_, tag, err := util.SplitImageRef(pilotImage)

	return tag, err
}
