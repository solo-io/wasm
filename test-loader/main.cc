#include <fstream>
#include <sstream>

#include "extensions/common/wasm/wasm.h"
#include "common/stats/isolated_store_impl.h"
#include "test/mocks/event/mocks.h"
#include "test/mocks/upstream/mocks.h"

int main() {

  NiceMock<Envoy::Event::MockDispatcher> dispatcher;
  NiceMock<Envoy::Upstream::MockClusterManager> cluster_manager;
  Envoy::Stats::IsolatedStoreImpl stats_store;
  std::ifstream t("../example/bazel-bin/envoy_filter_http_wasm_example.wasm");
  std::stringstream buffer;
  buffer << t.rdbuf();
  Envoy::Extensions::Common::Wasm::Wasm vm("envoy.wasm.runtime.v8", "vm_id", "vm_configuration",
           stats_store.createScope("wasm."), cluster_manager, dispatcher);
  if ( !vm.initialize(buffer.str(), true) ) {
    std::cerr << "error init" << std::endl;
    exit(1);
  }
}