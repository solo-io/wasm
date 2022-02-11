module github.com/solo-io/wasm/tools/wasme/pkg

go 1.15

require (
	github.com/avast/retry-go v2.4.3+incompatible
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/containerd/containerd v1.3.2
	github.com/deislabs/oras v0.8.1
	github.com/docker/distribution v2.8.0+incompatible
	github.com/docker/docker v1.13.1 // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.3.5
	github.com/google/go-cmp v0.5.0 // indirect
	github.com/gorilla/handlers v1.4.2 // indirect
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/hashicorp/go-multierror v1.0.0
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/miekg/dns v1.1.15 // indirect
	github.com/onsi/ginkgo v1.12.0
	github.com/onsi/gomega v1.8.1
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/opencontainers/image-spec v1.0.1
	github.com/opencontainers/runc v1.0.0-rc9 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/solo-io/go-utils v0.16.4
	github.com/stretchr/testify v1.6.1 // indirect
	go.opencensus.io v0.22.2 // indirect
	go.uber.org/zap v1.9.1
	golang.org/x/crypto v0.0.0-20200423211502-4bdfaf469ed5 // indirect
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b // indirect
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
	google.golang.org/genproto v0.0.0-20191115221424-83cc0476cb11 // indirect
	google.golang.org/grpc v1.28.1 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

replace (
	// copypaste from Gloo
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.0.0+incompatible
	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.4.2

	// Breaking changes pulled in by latest gloo need to use original repo instead of fork
	github.com/ilackarms/protoc-gen-doc => github.com/pseudomuto/protoc-gen-doc v1.3.0
)

replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309

replace github.com/codegangsta/cli => github.com/urfave/cli v1.22.2
