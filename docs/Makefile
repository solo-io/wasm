# Note to developers/testers
# >> to force docs publication and deployment during testing, do the following:
# cd to repo root, connect to a cluster where you want to push your images, run the following (update SUFFIX as you like)
# SUFFIX=a PROJECT_ID=solo-public TAGGED_VERSION=vtest-docs-build$SUFFIX make publish-docs -B && kubectl apply -f docs/install/manifest-latest.yaml -n docs
# >> to run host the docs locally, do the following
# cd to the docs dir
# make serve-site -B


# remove the "v" prefix
VERSION ?= $(shell echo $(TAGGED_VERSION) | cut -c 2-)

GCLOUD_PROJECT_ID ?= "solo-public"

GCR_REPO_PREFIX ?= gcr.io/$(GCLOUD_PROJECT_ID)
### REPLACE with product name
IMAGE_LEAF_NAME := web-assembly-hub-docs
IMAGE_REPO := $(GCR_REPO_PREFIX)/$(IMAGE_LEAF_NAME)

### REPLACE with product name
PRODUCT_SCOPE := web-assembly-hub
DOCS_VERSION := latest

#----------------------------------------------------------------------------------
# Docs
#----------------------------------------------------------------------------------

.PHONY: site-common
site-common: clean
# this theme is crucial (has the nested scoped short codes: protobuf and versioned_link_path - see web-assembly-hub docs for use demos)
	if [ ! -d themes/hugo-theme-soloio ]; then git clone https://github.com/solo-io/hugo-theme-soloio themes/hugo-theme-soloio; fi
	# style updates for putting docs in the web-assembly-hub repo, see details here https://github.com/solo-io/hugo-theme-soloio/commit/e0c50784a92fb7f61c635ff9a6e3a010f636f550
	git -C themes/hugo-theme-soloio checkout a9c18a63d56bea026a9e241ce0078caf56eabbc5

.PHONY: site-test
site-test: site-common
	hugo --config docs.toml --themesDir themes

.PHONY: site-release
site-release: site-common
	HUGO_PARAMS_Version=$(VERSION) hugo --config docs.toml --themesDir themes
	mv site site-latest

## Old way of deploying site. TODO: remove when old hosting setup is disabled
# .PHONY: deploy-site
# deploy-site: site
# 	firebase deploy --only hosting:gloo-docs

.PHONY: serve-site
serve-site: site-test
	HUGO_PARAMS_Version=$(VERSION) hugo --config docs.toml --themesDir themes server -D

.PHONY: clean
clean:
	rm -fr ./site ./resources ./site-latest ./site-versioned

# Uses https://github.com/gjtorikian/html-proofer
# Does not require running site; just make sure you generate the site and then run it
# Install with gem install html-proofer
# Another option we could use is wget: https://www.digitalocean.com/community/tutorials/how-to-find-broken-links-on-your-website-using-wget-on-debian-7
### rec: don't bother with this at the moment, requires some setup - not worth the effort for new docs
.PHONY: check-links
check-links:
	htmlproofer ./site/ --empty-alt-ignore  --allow-hash-href --alt-ignore "/img/logo.png" --url-ignore "/localhost/,/github.com/solo-io/solo-projects/,/developers.google.com/,/getgrav.org/,/github.com/solo-io/gloo/projects/,/developer.mozilla.org/"

# If on fedora, run
#     sudo dnf -y install gcc ruby-devel rubygems zlib-devel
# to install html-proofer deps (only works with gcc, not clang!)
install-tools:
	gem install html-proofer


# use vx.y.z-latest as the tag for images that are serving version x.y.z of the docs under the docs.solo.io/gloo/latest/ path
SERVE_AS_LATEST_TAG:=$(VERSION)-latest

# Let's say you only care about publishing to docs.solo.io/<product>/latest/
#  - I will comment out the simplifications
.PHONY: docker-push-docs
## don't bother generating manifest, just copy/modify from docs-prod repo
# docker-push-docs: site-release manifest
docker-push-docs: site-release
	@echo "${GCR_PASSWORD}" | docker login -u ${GCR_USERNAME} --password-stdin gcr.io
	docker build \
		--build-arg VERSION=latest \
		--build-arg PRODUCT_SCOPE=$(PRODUCT_SCOPE) \
		--build-arg FROM_DIR=./site-latest \
		-t $(GCR_REPO_PREFIX)/$(IMAGE_LEAF_NAME):$(SERVE_AS_LATEST_TAG) .
	docker push $(GCR_REPO_PREFIX)/$(IMAGE_LEAF_NAME):$(SERVE_AS_LATEST_TAG)
