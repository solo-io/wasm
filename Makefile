#----------------------------------------------------------------------------------
# Base
#----------------------------------------------------------------------------------

OUTDIR?=_output
PROJECT?=wasme

BUILDER_IMAGE?=quay.io/solo-io/ee-builder
CACHE_IMAGE?=quay.io/solo-io/wasme

SOURCES := $(shell find . -name "*.go" | grep -v test.go | grep -v '\.\#*')
RELEASE := "true"
ifeq ($(TAGGED_VERSION),)
	# TAGGED_VERSION := $(shell git describe --tags)
	# This doesn't work in CI, need to find another way...
	TAGGED_VERSION := vdev
	RELEASE := "false"
endif
VERSION ?= $(shell echo $(TAGGED_VERSION) | cut -c 2-)

LDFLAGS := "-X github.com/solo-io/$(PROJECT)/pkg/version.Version=$(VERSION)"
GCFLAGS := all="-N -l"

# Passed by cloudbuild
GCLOUD_PROJECT_ID := $(GCLOUD_PROJECT_ID)
BUILD_ID := $(BUILD_ID)

#----------------------------------------------------------------------------------
# Build
#----------------------------------------------------------------------------------

# must be a seperate target so that make waits for it to complete before moving on
.PHONY: mod-download
mod-download:
	go mod download

.PHONY: generate-deps
generate-deps: mod-download

# Build dependencies
.PHONY: generate-deps
generate-deps: mod-download
	go get -u github.com/cratonica/2goarray
	go get -u github.com/gogo/protobuf
	go get -u github.com/solo-io/protoc-gen-ext

# Generated Static assets for CLI & Docs, plus Operator/API Code
.PHONY: generated-code
generated-code: operator-gen
	go generate ./...

# Generate Operator Code & Chart
.PHONY: operator-gen
operator-gen:
	go run -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) operator/generate.go

# Generate Manifests from Chart
.PHONY: manifest-gen
manifest-gen: operator/install/kube/wasme-demo.yaml
operator/install/kube/wasme-demo.yaml: operator-gen
	helm template --namespace wasme operator/install/kube/wasme > operator/install/kube/wasme-demo.yaml

.PHONY: wasme
wasme: $(OUTDIR)/wasme
$(OUTDIR)/wasme: $(SOURCES)
	go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ cmd/main.go

.PHONY: wasme-linux-amd64
wasme-linux-amd64: $(OUTDIR)/wasme-linux-amd64
$(OUTDIR)/wasme-linux-amd64: $(SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ cmd/main.go

.PHONY: wasme-darwin-amd64
wasme-darwin-amd64: $(OUTDIR)/wasme-darwin-amd64
$(OUTDIR)/wasme-darwin-amd64: $(SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ cmd/main.go

.PHONY: wasme-windows-amd64
wasme-windows-amd64: $(OUTDIR)/wasme-windows-amd64.exe
$(OUTDIR)/wasme-windows-amd64.exe: $(SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ cmd/main.go


.PHONY: build-cli
build-cli: wasme-linux-amd64 wasme-darwin-amd64 wasme-windows-amd64

.PHONY: install-cli
install-cli:
	go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o ${GOPATH}/bin/wasme cmd/main.go


.PHONY: build-images
build-images: builder-image

# build Builder image
.PHONY: builder-image
builder-image:
	cd builder && \
	docker build -t $(BUILDER_IMAGE):$(VERSION) .

.PHONY: builder-image-push
builder-image-push:
	docker push $(BUILDER_IMAGE):$(VERSION)

#----------------------------------------------------------------------------------
# Release
#----------------------------------------------------------------------------------

# The code does the proper checking for a TAGGED_VERSION
.PHONY: upload-github-release-assets
upload-github-release-assets: build-cli
	go run ci/upload_github_release_assets.go

.PHONY: publish-docs
publish-docs:
ifeq ($(RELEASE),"true")
	$(MAKE) -C docs docker-push-docs \
		VERSION=$(VERSION) \
		TAGGED_VERSION=$(TAGGED_VERSION) \
		GCLOUD_PROJECT_ID=$(GCLOUD_PROJECT_ID) \
		RELEASE=$(RELEASE)
endif

.PHONY: publish-images
publish-images:
ifeq ($(RELEASE),"true")
	docker push $(BUILDER_IMAGE):$(VERSION)
	docker push $(CACHE_IMAGE):$(VERSION)
endif

#----------------------------------------------------------------------------------
# Clean
#----------------------------------------------------------------------------------

# Important to clean before pushing new releases. Dockerfiles and binaries may not update properly
.PHONY: clean
clean:
	rm -rf  _output/
	rm -rf  example/cpp/{bazel-bin,bazel-out,bazel-testlogs,bazel-workspace}

.PHONY: build-example
build-example:
	go run cmd/main.go build example/cpp
