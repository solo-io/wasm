.PHONY: image
image:
	docker build -t soloio/wasm-builder .

.PHONY: build
build:
	./build.sh example/cpp _output

.PHONY: clean
clean:
	rm -rf  ./build_output && \
	rm -rf ./workspace && \

