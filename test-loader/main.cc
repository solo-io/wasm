#include <fstream>
#include <sstream>

#include "extensions/common/wasm/wasm.h"
#include "common/stats/isolated_store_impl.h"
#include "test/mocks/event/mocks.h"
#include "test/mocks/upstream/mocks.h"

int main(int argc, char* argv[]) {

  if (args != 2) {
    std::cerr << "please provide path to wasm file" << std::endl;
    exit(1);
  }

  const char* wasm_module = argv[1];

  NiceMock<Envoy::Event::MockDispatcher> dispatcher;
  NiceMock<Envoy::Upstream::MockClusterManager> cluster_manager;
  Envoy::Stats::IsolatedStoreImpl stats_store;
  std::ifstream t(wasm_module);
  std::stringstream buffer;
  buffer << t.rdbuf();
  Envoy::Extensions::Common::Wasm::Wasm vm("envoy.wasm.runtime.v8", "vm_id", "vm_configuration",
           stats_store.createScope("wasm."), cluster_manager, dispatcher);
  if ( !vm.initialize(buffer.str(), true) ) {
    std::cerr << "error init" << std::endl;
    exit(1);
  }
  std::cout << "loaded successfully" << std::endl;
  return 0;
}