package local_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/solo-io/wasme/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/wasme/pkg/defaults"
	v1 "github.com/solo-io/wasme/pkg/operator/api/wasme.io/v1"
	"github.com/solo-io/wasme/pkg/store"
	"github.com/solo-io/wasme/pkg/util"

	. "github.com/solo-io/wasme/pkg/deploy/local"
)

var _ = Describe("LocalProvider", func() {
	var filter = &v1.FilterSpec{
		Id:     "my_filter",
		Image:  "yuvaltest.solo.io/ilackarms/assemblyscript-test:istio-1.5.0-alpha.0",
		Config: "wurld",
	}
	It("prints the injected yaml", func() {

		err := test.WasmeCli("pull", filter.Image)
		Expect(err).NotTo(HaveOccurred())

		store := store.NewStore(defaults.WasmeImageDir)

		buf := &bytes.Buffer{}
		p := &Runner{
			Ctx:              context.TODO(),
			Input:            ioutil.NopCloser(bytes.NewBuffer([]byte(BasicEnvoyConfig))),
			Output:           buf,
			Store:            store,
			DockerRunArgs:    nil,
			EnvoyArgs:        nil,
			EnvoyDockerImage: "",
		}
		err = p.RunFilter(filter)
		Expect(err).NotTo(HaveOccurred())

		filterDir, err :=store.Dir(filter.Image)
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
			Store:            store.NewStore(defaults.WasmeImageDir),
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

		// test with curl!

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

			return out, err

			//log.Printf("output: %v", out)
			//log.Printf("err: %v", err)
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
