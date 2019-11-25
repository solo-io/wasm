licenses(["notice"])  # Apache 2

load(
    "@envoy//bazel:envoy_build_system.bzl",
    "envoy_cc_binary",
    "envoy_package",
)

envoy_package()

envoy_cc_binary(
    name = "testloader",
    repository = "@envoy",
    testonly = True,
    srcs = ["main.cc"],
    deps = [
        "@envoy//source/extensions/common/wasm:wasm_lib",
        "@envoy//source/common/stats:isolated_store_lib",
        "@envoy//test/mocks/upstream:upstream_mocks",
        "@envoy//test/mocks/event:event_mocks",
    ],
)
