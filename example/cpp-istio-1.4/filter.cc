// NOLINT(namespace-envoy)
#include <string>
#include <unordered_map>

#include "google/protobuf/util/json_util.h"
#include "proxy_wasm_intrinsics.h"
#include "filter.pb.h"

class AddHeaderRootContext : public RootContext {
public:
  explicit AddHeaderRootContext(uint32_t id, StringView root_id) : RootContext(id, root_id) {}
  bool onConfigure(std::unique_ptr<WasmData> conf) override;

  std::string header_value_;
};

class AddHeaderContext : public Context {
public:
  explicit AddHeaderContext(uint32_t id, RootContext* root) : Context(id, root), root_(static_cast<AddHeaderRootContext*>(static_cast<void*>(root))) {}

  void onCreate() override;
  FilterHeadersStatus onResponseHeaders() override;
private:

  AddHeaderRootContext* root_;
};

bool AddHeaderRootContext::onConfigure(std::unique_ptr<WasmData> conf) { 
  Config config;
  
  google::protobuf::util::JsonParseOptions options;
  options.case_insensitive_enum_parsing = true;
  options.ignore_unknown_fields = false;

  google::protobuf::util::JsonStringToMessage(conf->toString(), &config, options);
  LOG_DEBUG("onConfigure " + config.value());
  header_value_ = config.value();
  return true; 
}

void AddHeaderContext::onCreate() { LOG_DEBUG(std::string("onCreate " + std::to_string(id()))); }

FilterHeadersStatus AddHeaderContext::onResponseHeaders() {
  LOG_DEBUG(std::string("onResponseHeaders ") + std::to_string(id()));
  addResponseHeader("newheader", root_->header_value_);
  return FilterHeadersStatus::Continue;
}

static RegisterContextFactory register_AddHeaderContext(CONTEXT_FACTORY(AddHeaderContext),
                                                      ROOT_FACTORY(AddHeaderRootContext),
                                                      "add_header_root_id");