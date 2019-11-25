// NOLINT(namespace-envoy)
#include <string>
#include <unordered_map>

#include "proxy_wasm_intrinsics.h"

class ExampleRootContext : public RootContext {
public:
  explicit ExampleRootContext(uint32_t id, StringView root_id) : RootContext(id, root_id) {}
  bool onConfigure(size_t /* configuration_size */) override;

  bool onStart(size_t) override;

  std::string yuval_;
};

class ExampleContext : public Context {
public:
  explicit ExampleContext(uint32_t id, RootContext* root) : Context(id, root), root_(static_cast<ExampleRootContext*>(static_cast<void*>(root))) {}

  void onCreate() override;
  FilterHeadersStatus onRequestHeaders(uint32_t headers) override;
  FilterDataStatus onRequestBody(size_t body_buffer_length, bool end_of_stream) override;
  FilterHeadersStatus onResponseHeaders(uint32_t headers) override;
  void onDone() override;
  void onLog() override;
  void onDelete() override;

  ExampleRootContext* root_;
};
static RegisterContextFactory register_ExampleContext(CONTEXT_FACTORY(ExampleContext),
                                                      ROOT_FACTORY(ExampleRootContext),
                                                      "my_root_id");

bool ExampleRootContext::onConfigure(size_t) { 
  auto conf = getConfiguration();
  LOG_WARN("onConfigure1 " + std::to_string(conf->size()));
  LOG_WARN("onConfigure2 " + conf->toString());
  yuval_ = conf->toString();
  if (yuval_ != "yuval") {
    return true;
  }
  return true; 
}

bool ExampleRootContext::onStart(size_t) { LOG_WARN("onStart"); return true;}

void ExampleContext::onCreate() { LOG_WARN(std::string("onCreate " + std::to_string(id()))); }

FilterHeadersStatus ExampleContext::onRequestHeaders(uint32_t) {
  LOG_INFO(std::string("onRequestHeaders ") + std::to_string(id()));
  LOG_INFO(std::string("onRequestHeaders ") + root_->yuval_);
  auto result = getRequestHeaderPairs();
  auto pairs = result->pairs();
  LOG_INFO(std::string("headers: ") + std::to_string(pairs.size()));
  for (auto& p : pairs) {
    LOG_INFO(std::string(p.first) + std::string(" -> ") + std::string(p.second));
  }
  return FilterHeadersStatus::Continue;
}

FilterHeadersStatus ExampleContext::onResponseHeaders(uint32_t) {
  LOG_DEBUG(std::string("onResponseHeaders ") + std::to_string(id()));
  auto result = getResponseHeaderPairs();
  auto pairs = result->pairs();
  LOG_INFO(std::string("headers: ") + std::to_string(pairs.size()));
  for (auto& p : pairs) {
    LOG_INFO(std::string(p.first) + std::string(" -> ") + std::string(p.second));
  }
  addResponseHeader("newheader", "joe");
  replaceResponseHeader("location", "envoy-wasm");
  return FilterHeadersStatus::Continue;
}

FilterDataStatus ExampleContext::onRequestBody(size_t body_buffer_length, bool end_of_stream) {
  auto body = getBufferBytes(BufferType::HttpRequestBody, 0, body_buffer_length);
  LOG_ERROR(std::string("onRequestBody ") + std::string(body->view()));
  return FilterDataStatus::Continue;
}

void ExampleContext::onDone() { LOG_WARN(std::string("onDone " + std::to_string(id()))); }

void ExampleContext::onLog() { LOG_WARN(std::string("onLog " + std::to_string(id()))); }

void ExampleContext::onDelete() { LOG_WARN(std::string("onDelete " + std::to_string(id()))); }
