module github.com/solo-io/wasme

go 1.13

require (
	github.com/containerd/containerd v1.3.0
	github.com/deislabs/oras v0.7.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/google/go-github/v28 v28.1.1
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/opencontainers/image-spec v1.0.1
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/solo-io/gloo v1.2.10
	github.com/solo-io/go-utils v0.11.0
	github.com/spf13/cobra v0.0.5
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	gopkg.in/AlecAivazis/survey.v1 v1.8.7
	istio.io/client-go v0.0.0-20191206191348-5c576a7ecef0
	k8s.io/api v0.0.0
)

// Pinned to kubernetes-1.14.1
replace (
	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.4.2
	k8s.io/api => k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190409022649-727a075fdec8
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20190409021813-1ec86e4da56c
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20190409023024-d644b00f3b79
	k8s.io/client-go => k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20190409023720-1bc0c81fa51d
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20190311093542-50b561225d70
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20190409022021-00b8e31abe9d
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190510232812-a01b7d5d6c22
	k8s.io/kubectl => k8s.io/kubectl v0.0.0-20191016120415-2ed914427d51
	k8s.io/kubernetes => k8s.io/kubernetes v1.14.1
)

replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309

replace github.com/codegangsta/cli => github.com/urfave/cli v1.22.2
