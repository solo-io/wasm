module github.com/solo-io/extend-envoy

go 1.13

require (
	github.com/containerd/containerd v1.3.0
	github.com/deislabs/oras v0.7.0
	github.com/golang/protobuf v1.3.2
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/opencontainers/image-spec v1.0.1
	github.com/spf13/cobra v0.0.5
	k8s.io/apimachinery v0.0.0-20191111054156-6eb29fdf75dc
)

replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309
