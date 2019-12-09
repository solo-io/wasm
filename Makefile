IMAGE?=quay.io/solo-io/ee-builder
TAG?=v1

# install codegen deps
.PHONY: gen-deps
gen-deps:
	go get -u github.com/cratonica/2goarray

# generate code (static assets for CLI)
.PHONY: generated-code
generated-code:
	go generate ./...

# build Builder image
.PHONY: image
image:
	cd builder && \
	docker build -t $(IMAGE):$(TAG) .

.PHONY: image-push
image-push:
	docker push $(IMAGE):$(TAG)

.PHONY: build-example
build-example:
	go run main.go build example/cpp

.PHONY: clean
clean:
	rm -rf  _output/
	rm -rf  example/cpp/{bazel-bin,bazel-out,bazel-testlogs,bazel-workspace}
