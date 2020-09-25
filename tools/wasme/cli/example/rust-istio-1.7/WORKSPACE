workspace(name = "filter_example")

load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository", "new_git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")


http_archive(
    name = "proxy_wasm_rust_sdk",
    sha256 = "3b1b8912a322b80e101686be61139000994594f4fad6bf7c5546b1d35d9cb8c2",
    strip_prefix = "proxy-wasm-rust-sdk-e7a6ad0f123965b5c05fa85606da9da00cac6ab6",
    urls = [
        "https://github.com/proxy-wasm/proxy-wasm-rust-sdk/archive/e7a6ad0f123965b5c05fa85606da9da00cac6ab6.tar.gz",
    ],
)

http_archive(
    name = "bazel_skylib",
    sha256 = "97e70364e9249702246c0e9444bccdc4b847bed1eb03c5a3ece4f83dfe6abc44",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.0.2/bazel-skylib-1.0.2.tar.gz",
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.0.2/bazel-skylib-1.0.2.tar.gz",
    ],
)

load("@bazel_skylib//:workspace.bzl", "bazel_skylib_workspace")

bazel_skylib_workspace()

http_archive(
    name = "io_bazel_rules_rust",
    sha256 = "484a2b2b67cd2d1fa1054876de7f8d291c4b203fd256bc8cbea14d749bb864ce",
    # Last commit where "out_binary = True" works.
    # See: https://github.com/bazelbuild/rules_rust/issues/386
    strip_prefix = "rules_rust-fda9a1ce6482973adfda022cadbfa6b300e269c3",
    url = "https://github.com/bazelbuild/rules_rust/archive/fda9a1ce6482973adfda022cadbfa6b300e269c3.tar.gz",
)

load("@io_bazel_rules_rust//rust:repositories.bzl", "rust_repositories")

rust_repositories()

load("@io_bazel_rules_rust//:workspace.bzl", "bazel_version")

bazel_version(name = "bazel_version")

load("@proxy_wasm_rust_sdk//bazel/cargo:crates.bzl", "raze_fetch_remote_crates")

raze_fetch_remote_crates()
