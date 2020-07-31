package local_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/solo-io/wasme/pkg/consts"
	"github.com/solo-io/wasme/pkg/store"
	"github.com/solo-io/wasme/pkg/util"
	"github.com/solo-io/wasme/test"

	"github.com/gogo/protobuf/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/wasme/pkg/deploy/local"

	testvars "github.com/solo-io/wasme/pkg/consts/test"
	v1 "github.com/solo-io/wasme/pkg/operator/api/wasme.io/v1"
)

var _ = Describe("LocalProvider", func() {
	sv := &types.StringValue{
		Value: "wurld",
	}
	val, _ := sv.Marshal()
	var (
		filter = &v1.FilterSpec{
			Id:    "my_filter",
			Image: consts.HubDomain + "/ilackarms/assemblyscript-test:" + testvars.Istio15Tag,
			Config: &types.Any{
				TypeUrl: "type.googleapis.com/google.protobuf.StringValue",
				Value:   val,
			},
		}
		storeDir string

		imageStore store.Store
	)
	BeforeEach(func() {
		// need to run with storage dir set to . in CI due to docker mount concerns
		dir, err := ioutil.TempDir(".", "local-wasme-test")
		Expect(err).NotTo(HaveOccurred())
		storeDir = dir

		err = test.WasmeCli("pull", filter.Image, "--store="+storeDir)
		Expect(err).NotTo(HaveOccurred())

		imageStore = store.NewStore(storeDir)

	})
	AfterEach(func() {
		os.RemoveAll(storeDir)
	})
	It("runs Envoy locally with the filter mounted", func() {
		buf := &bytes.Buffer{}
		p := &Runner{
			Ctx:              context.TODO(),
			Input:            ioutil.NopCloser(bytes.NewBuffer([]byte(BasicEnvoyConfig))),
			Output:           buf,
			Store:            imageStore,
			DockerRunArgs:    []string{"--net=host"},
			EnvoyArgs:        nil,
			EnvoyDockerImage: "",
		}
		err := p.RunFilter(filter)
		Expect(err).NotTo(HaveOccurred())

		filterDir, err := imageStore.Dir(filter.Image)
		Expect(err).NotTo(HaveOccurred())

		Expect(buf.String()).To(Equal(expectedConfig(filterDir)))
	})
	AfterEach(func() {
		util.Docker(nil, nil, nil, "kill", filter.Id)
	})
	It("runs envoy locally with the given filter", func() {

		p := &Runner{
			Ctx:              context.TODO(),
			Input:            ioutil.NopCloser(bytes.NewBuffer([]byte(BasicEnvoyConfig))),
			Store:            imageStore,
			DockerRunArgs:    nil,
			EnvoyArgs:        nil,
			EnvoyDockerImage: DefaultEnvoyImage,
		}
		var runError error
		var errLock sync.RWMutex
		go func() {
			err := p.RunFilter(filter)
			errLock.Lock()
			runError = err
			errLock.Unlock()
		}()

		// envoy addr defaults to the localhost port-forward

		// test with curl!
		t := time.Tick(time.Second)

		testRequest := func() (string, error) {
			errLock.RLock()
			Expect(runError).NotTo(HaveOccurred())
			errLock.RUnlock()
			b := &bytes.Buffer{}
			err := util.ExecCmd(
				b,
				b,
				nil,
				"curl",
				"-v",
				"localhost:8080/")

			out := b.String()
			select {
			case <-t:
				log.Printf("out: %v", out)
				log.Printf("err: %v", err)
			default:
			}

			return out, err
		}

		// expect header in response
		Eventually(testRequest, time.Minute*5).Should(ContainSubstring("hello: wurld"))
	})
})

func expectedConfig(dir string) string {
	return fmt.Sprintf(`admin:
  accessLogPath: /dev/null
  address:
    socketAddress:
      address: 0.0.0.0
      portValue: 19000
staticResources:
  clusters:
  - connectTimeout: 0.250s
    dnsLookupFamily: V4_ONLY
    hosts:
    - socketAddress:
        address: jsonplaceholder.typicode.com
        ipv4Compat: true
        portValue: 443
    name: static-cluster
    tlsContext:
      sni: jsonplaceholder.typicode.com
    type: LOGICAL_DNS
  listeners:
  - address:
      socketAddress:
        address: 0.0.0.0
        portValue: 8080
    filterChains:
    - filters:
      - config:
          httpFilters:
          - config:
              config:
                configuration: wurld
                name: my_filter
                rootId: add_header
                vmConfig:
                  code:
                    local:
                      filename: %v/filter.wasm
                  runtime: envoy.wasm.runtime.v8
                  vmId: my_filter
            name: envoy.filters.http.wasm
          - name: envoy.router
          routeConfig:
            name: test
            virtualHosts:
            - domains:
              - '*'
              name: jsonplaceholder
              routes:
              - match:
                  prefix: /
                route:
                  autoHostRewrite: true
                  cluster: static-cluster
          statPrefix: ingress_http
        name: envoy.http_connection_manager
    name: listener_0
`, dir)
}
