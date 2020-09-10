This compiles an example filter for envoy WASM.

# build filter
build with
```
bazel build :filter.wasm
```

Filter will be in:
```
./bazel-bin/filter.wasm
```

# build config descriptors

build descriptors with:
```
bazel build :filter_proto
```

Descriptors will be in:
```
./bazel-bin/filter_proto-descriptor-set.proto.bin
```

Note: 
on a mac, please run
```
xcode-select --install
```

and Potentially:
```
brew install python@2
```
as the python bundled with catalina may have issues with ssl certs.