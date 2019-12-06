.PHONY: image
image:
	cd builder && \
	docker build -t soloio/wasm-builder .

.PHONY: build-example
build-example:
	builder/build.sh $(PWD)/example/cpp $(PWD)/_output

.PHONY: clean
clean:
	rm -rf  ./builder/build_output && \
	rm -rf ./builder/workspace && \

