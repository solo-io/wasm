
cc_binary(
    name = "filter.wasm",
    srcs = [
        "filter.cc"
    ],
    deps = [
        "@envoy_wasm_api//:proxy_wasm_intrinsics",
    ],
    additional_linker_inputs = ["@envoy_wasm_api//:jslib"],
    linkopts = ["--js-library", "external/envoy_wasm_api/proxy_wasm_intrinsics.js"],
)