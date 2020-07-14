package(default_visibility = ["//visibility:public"])

load(":cc_toolchain_config.bzl", "cc_toolchain_config")

cc_toolchain_suite(
    name = "emscripten",
    toolchains = {
        "wasm": ":wasm_toolchain",
    },
)

filegroup(name = "empty")
 
filegroup(
  name = "emscripten_cache_content",
  srcs = glob(["tmp/emscripten_cache/**/*"]),
)

filegroup(
    name = "all",
    srcs = [
        "emcc.sh",
        "emar.sh",
        "common.sh",
        "@emscripten_toolchain//:all",
        ":emscripten_cache_content"
    ],
)
cc_toolchain(
    name = "wasm_toolchain",
    toolchain_identifier = "wasm-toolchain",
    toolchain_config = ":wasm_toolchain_config",
    all_files = ":all",
    compiler_files = ":all",
    ar_files = ":all",
    dwp_files = ":empty",
    linker_files = ":all",
    objcopy_files = ":empty",
    strip_files = ":empty",
    supports_param_files = 0,
)

cc_toolchain_config(name = "wasm_toolchain_config")
