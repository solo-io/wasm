#----------------------------------------------------------------------------------
# Base
#----------------------------------------------------------------------------------

OUTDIR?=_output
PROJECT?=wasme

BUILDER_IMAGE?=quay.io/solo-io/ee-builder

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

.PHONY: enable-gomod
enable-gomod:
	export GO11MODULE=on

# Build dependencies
.PHONY: generate-deps
generate-deps: enable-gomod
	go get -u github.com/cratonica/2goarray

# Generated Static assets for CLI & Docs
.PHONY: generated-code
generated-code: enable-gomod
	go generate ./...

.PHONY: wasme
wasme: $(OUTDIR)/wasme
$(OUTDIR)/wasme: $(SOURCES) enable-gomod
	go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ main.go

.PHONY: wasme-linux-amd64
wasme-linux-amd64: $(OUTDIR)/wasme-linux-amd64
$(OUTDIR)/wasme-linux-amd64: $(SOURCES) enable-gomod
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ main.go

.PHONY: wasme-darwin-amd64
wasme-darwin-amd64: $(OUTDIR)/wasme-darwin-amd64
$(OUTDIR)/wasme-darwin-amd64: $(SOURCES) enable-gomod
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ main.go

.PHONY: wasme-windows-amd64
wasme-windows-amd64: $(OUTDIR)/wasme-windows-amd64.exe
$(OUTDIR)/wasme-windows-amd64.exe: $(SOURCES) enable-gomod
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ main.go


.PHONY: build-cli
build-cli: wasme-linux-amd64 wasme-darwin-amd64 wasme-windows-amd64

.PHONY: install-cli
install-cli: enable-gomod
	go build -o ${GOPATH}/bin/wasme main.go


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
	go run main.go build example/cpp
