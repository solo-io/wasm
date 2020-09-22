load("@rules_proto//proto:defs.bzl", "proto_library")
load("//bazel/wasm:wasm.bzl", "wasm_cc_binary")

proto_library(
    name = "filter_proto",
    srcs = [
        "filter.proto",
    ],
)

cc_proto_library(
    name = "filter_cc_proto",
    deps = [":filter_proto"],
)

wasm_cc_binary(
    name = "filter.wasm",
    srcs = [
        "filter.cc",
    ],
    additional_linker_inputs = ["@proxy_wasm_cpp_sdk//:jslib"],
    linkopts = [
        "--js-library",
        "external/proxy_wasm_cpp_sdk/proxy_wasm_intrinsics.js",
    ],
    deps = [
        ":filter_cc_proto",
        "@proxy_wasm_cpp_sdk//:proxy_wasm_intrinsics",
    ],
)
