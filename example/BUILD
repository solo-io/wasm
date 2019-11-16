
cc_binary(
    name = "envoy_filter_http_wasm_example.wasm",
    srcs = [
        "envoy_filter_http_wasm_example.cc"
    ],
    deps = [
        "@envoy_wasm//:proxy_wasm_intrinsics",
    ],
    additional_linker_inputs = ["@envoy_wasm//:jslib"],
    linkopts = ["--js-library", "external/envoy_wasm/proxy_wasm_intrinsics.js"],
)