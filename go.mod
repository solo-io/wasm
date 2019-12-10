module github.com/solo-io/wasme

go 1.13

require (
	github.com/containerd/containerd v1.3.0
	github.com/deislabs/oras v0.7.0
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/google/go-github/v28 v28.1.1
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/opencontainers/image-spec v1.0.1
	github.com/sirupsen/logrus v1.4.2
	github.com/solo-io/go-utils v0.11.0
	github.com/spf13/cobra v0.0.5
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sys v0.0.0-20191206220618-eeba5f6aabab // indirect
	gopkg.in/AlecAivazis/survey.v1 v1.8.7
)

replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309

replace github.com/codegangsta/cli => github.com/urfave/cli v1.22.2
