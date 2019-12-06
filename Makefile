.PHONY: image
image:
	cd builder && \
	docker build -t soloio/wasm-builder .

.PHONY: build-example
build-example:
	builder/build.sh $(PWD)/example/cpp $(PWD)/_output

.PHONY: clean
clean:
	rm -rf  _output/
	rm -rf  example/cpp/{bazel-bin,bazel-out,bazel-testlogs,bazel-workspace}


