IMAGE?=quay.io/solo-io/ee-builder
TAG?=v1

.PHONY: image
image:
	cd builder && \
	docker build -t $(IMAGE):$(TAG) .

.PHONY: build-example
build-example:
	go run main.go build example/cpp

.PHONY: clean
clean:
	rm -rf  _output/
	rm -rf  example/cpp/{bazel-bin,bazel-out,bazel-testlogs,bazel-workspace}
