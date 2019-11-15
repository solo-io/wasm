
# add proto.a in @envoy//api/wasm/cpp/libprotobuf-lite.a 
# add includes in @envoy//api/wasm/cpp/proxy_wasm_intrinsics.h,proxy_wasm_enums.h,proxy_wasm_externs.h,proxy_wasm_intrinsics_lite.pb.cc,struct_lite.pb.cc

a="""
envoy_cc_library(
    name = "protbuf",
    srcs = [
        "@envoy//api/wasm/cpp/proxy_wasm_intrinsics.h",
        "@envoy//api/wasm/cpp/proxy_wasm_enums.h",
        "@envoy//api/wasm/cpp/proxy_wasm_externs.h",
    ],
    hdrs = [
        "@envoy//api/wasm/cpp/proxy_wasm_intrinsics_lite.pb.cc",
        "@envoy//api/wasm/cpp/struct_lite.pb.cc",
    ],
    repository = "@envoy",
    deps = [
    ],
)
"""

cc_binary(
    name = "yuval",
    srcs = [
        "envoy_filter_http_wasm_example.cc"
    ],
    deps = [
        "@envoy_wasm//:proxy_wasm_intrinsics",
    ],
)