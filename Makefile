#----------------------------------------------------------------------------------
# Base
#----------------------------------------------------------------------------------

OUTDIR?=_output
PROJECT?=wasme

BUILDER_IMAGE?=quay.io/solo-io/ee-builder
OPERATOR_IMAGE?=quay.io/solo-io/wasme

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
GCLOUD_PROJECT_ID?=solo-public
BUILD_ID := $(BUILD_ID)

#----------------------------------------------------------------------------------
# Build
#----------------------------------------------------------------------------------

# Build dependencies

.PHONY: mod-download
mod-download:
	go mod download

.PHONY: install-deps
install-deps: mod-download
	go get -v istio.io/tools/cmd/protoc-gen-jsonshim@v0.0.0-20200414140130-90db7d74fac2
	go get -v github.com/gogo/protobuf/protoc-gen-gogo@v1.3.1
	# github.com/golang/protobuf pinned to 1.3.5 else docs/generate_reference_docs.go proto.Unmarshal panics
	go get -v github.com/golang/protobuf/protoc-gen-go@v1.3.5
	go get -v github.com/solo-io/protoc-gen-ext@v0.0.9
	go get -v github.com/golang/mock/mockgen@v1.4.3
	go get -v golang.org/x/tools/cmd/goimports@v0.0.0-20200414131530-0037cb7812fa
	go get -v github.com/cratonica/2goarray
	go get -v github.com/solo-io/gloo@v1.5.0-beta11
	go get -v github.com/solo-io/skv2@v0.8.0

	go mod tidy



# Generated Static assets for CLI & Docs, plus Operator/API Code
SUBDIRS:=pkg test operator
.PHONY: generated-code
generated-code:
	go generate ./...
	goimports -w $(SUBDIRS)

# Generate Operator Code & Chart
.PHONY: operator-gen
operator-gen:
	go run -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) operator/generate.go
	goimports -w $(SUBDIRS)

# Generate Manifests from Chart
.PHONY: manifest-gen
manifest-gen: operator/install/wasme-default.yaml
operator/install/wasme-default.yaml: operator-gen
	helm template --namespace wasme operator/install/wasme > operator/install/wasme-default.yaml

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
build-images: wasme-image builder-image

# build image with Wasme binary
.PHONY: wasme-image
wasme-image: wasme-linux-amd64
	cp $(OUTDIR)/wasme-linux-amd64 operator/build/wasme/ && \
	docker build -t $(OPERATOR_IMAGE):$(VERSION) operator/build/wasme/
	rm operator/build/wasme/wasme-linux-amd64

.PHONY: wasme-image-push
wasme-image-push:
	docker push $(OPERATOR_IMAGE):$(VERSION)

# build Builder image
.PHONY: builder-image
builder-image:
	cd builder && \
	docker build -t $(BUILDER_IMAGE):$(VERSION) .

.PHONY: builder-image-push
builder-image-push:
	docker push $(BUILDER_IMAGE):$(VERSION)

#----------------------------------------------------------------------------------
# Test
#----------------------------------------------------------------------------------

# run all tests
# set TEST_PKG to run a specific test package
.PHONY: run-tests
run-tests:
	ginkgo -v -r -failFast -trace -progress \
		-ldflags=$(LDFLAGS) \
		-gcflags=$(GCFLAGS) \
		-progress \
		-compilers=4 \
		-skipPackage=$(SKIP_PACKAGES) $(TEST_PKG)

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
	make wasme-image-push
	make builder-image-push
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
