module github.com/solo-io/wasme

go 1.13

require (
	github.com/avast/retry-go v2.4.3+incompatible
	github.com/containerd/containerd v1.3.2
	github.com/deislabs/oras v0.8.1
	github.com/docker/cli v0.0.0-20200130152716-5d0cf8839492
	github.com/docker/distribution v2.7.1+incompatible
	github.com/envoyproxy/go-control-plane v0.9.6-0.20200529035633-fc42e08917e9
	github.com/envoyproxy/protoc-gen-validate v0.4.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/gogo/protobuf v1.3.1
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.2
	github.com/google/go-cmp v0.5.0 // indirect
	github.com/hashicorp/go-multierror v1.0.0
	github.com/manifoldco/promptui v0.7.0
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/opencontainers/image-spec v1.0.1
	github.com/pkg/errors v0.9.1
	github.com/pseudomuto/protoc-gen-doc v1.3.2
	github.com/pseudomuto/protokit v0.2.0
	github.com/sirupsen/logrus v1.6.0
	github.com/solo-io/gloo v1.5.0-beta11
	github.com/solo-io/go-utils v0.17.0
	github.com/solo-io/protoc-gen-ext v0.0.9
	github.com/solo-io/skv2 v0.8.0
	github.com/solo-io/solo-kit v0.13.11
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	go.uber.org/zap v1.15.0
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	istio.io/api v0.0.0-20191109011911-e51134872853
	istio.io/client-go v0.0.0-20191206191348-5c576a7ecef0
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/code-generator v0.18.6
	sigs.k8s.io/controller-runtime v0.6.2
)

// Pinned to kubernetes-1.18
replace (
	// copypaste from Gloo
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.0.0+incompatible
	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.4.2

	// Breaking changes pulled in by latest gloo need to use original repo instead of fork
	github.com/ilackarms/protoc-gen-doc => github.com/pseudomuto/protoc-gen-doc v1.3.0
	github.com/solo-io/gloo => github.com/solo-io/gloo v1.5.0-beta20

	github.com/solo-io/go-utils => github.com/solo-io/go-utils v0.17.0
	github.com/solo-io/solo-kit => github.com/solo-io/solo-kit v0.14.0

	k8s.io/api => k8s.io/api v0.18.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.6
	k8s.io/apiserver => k8s.io/apiserver v0.18.6
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.18.6
	k8s.io/client-go => k8s.io/client-go v0.18.6
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.18.6
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.18.6
	k8s.io/code-generator => k8s.io/code-generator v0.18.6
	k8s.io/component-base => k8s.io/component-base v0.18.6
	k8s.io/cri-api => k8s.io/cri-api v0.18.6
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.18.6
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.18.6
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.18.6
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.18.6
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.18.6
	k8s.io/kubectl => k8s.io/kubectl v0.18.6
	k8s.io/kubelet => k8s.io/kubelet v0.18.6
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.18.6
	k8s.io/metrics => k8s.io/metrics v0.18.6
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.18.6
	k8s.io/utils => k8s.io/utils v0.0.0-20200821003339-5e75c0163111
)

replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309

replace github.com/codegangsta/cli => github.com/urfave/cli v1.22.2
