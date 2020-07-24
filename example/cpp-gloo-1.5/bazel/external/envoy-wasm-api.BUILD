package(default_visibility = ['//visibility:public'])

cc_library(
    name = "proxy_wasm_intrinsics",
    visibility = ["//visibility:public"],
    srcs = [
        "proxy_wasm_intrinsics.cc",
        "proxy_wasm_intrinsics_lite.pb.cc",
        "struct_lite.pb.cc",
    ],
    hdrs = [
        "proxy_wasm_intrinsics.h",
        "proxy_wasm_enums.h",
        "proxy_wasm_common.h",
        "proxy_wasm_externs.h",
        "proxy_wasm_api.h",
        "proxy_wasm_intrinsics.pb.h",
        "proxy_wasm_intrinsics_lite.pb.h",
        "struct_lite.pb.h",
    ],
    deps = [
        "@com_google_protobuf//:protobuf_lite",
    ],
)
filegroup(
    visibility = ["//visibility:public"],
    name = "jslib",
    srcs = [
        "proxy_wasm_intrinsics.js",
    ],
)