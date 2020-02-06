# Note to developers/testers
# >> to force docs publication and deployment during testing, do the following:
# cd to repo root, connect to a cluster where you want to push your images, run the following (update SUFFIX as you like)
# SUFFIX=a PROJECT_ID=solo-public TAGGED_VERSION=vtest-docs-build$SUFFIX make publish-docs -B && kubectl apply -f docs/install/manifest-latest.yaml -n docs
# >> to run host the docs locally, do the following
# cd to the docs dir
# make serve-site -B



# remove the "v" prefix
VERSION ?= $(shell echo $(TAGGED_VERSION) | cut -c 2-)

# GCLOUD_PROJECT_ID Passed by cloudbuild upload-release-assets
GCLOUD_PROJECT_ID := $(GCLOUD_PROJECT_ID)

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
	git -C themes/hugo-theme-soloio checkout e0c50784a92fb7f61c635ff9a6e3a010f636f550

.PHONY: site-test
site-test: site-common
### skip, to host locally
	# go run cmd/generate_changelog_doc.go gen-version-scope-data --no-scope
## added this as an alternative to the build script:
	cp data/Solo.testlocal.yaml data/Solo.yaml
	hugo --config docs.toml
	# ensure that valid json is generated. Common cause: using yaml ">" multiline string symbols in Hugo's toml headers
	cat site/index.json | jq "." > /dev/null

.PHONY: site-release
site-release: site-common
### This generates docs/data/Solo.yaml which indicates the docs (whether "latest" or x.y.z) and code version (x.y.z)
### For simplicity, I have hard coded that file here - just set the product/version as you like, keep the docs version at /<product>/latest
	# go run cmd/generate_changelog_doc.go gen-version-scope-data --product gloo --version $(VERSION) --call-latest
## added this as an alternative to the build script:
	cp data/Solo.publish.yaml data/Solo.yaml
	hugo --config docs.toml
	# ensure that valid json is generated. Common cause: using yaml ">" multiline string symbols in Hugo's toml headers
	# (if it passes here, it will pass on the subsequent generation so no need to check both hugo calls)
	cat site/index.json | jq "." > /dev/null
	mv site site-latest
### Skip versioned release
	# go run cmd/generate_changelog_doc.go gen-version-scope-data --product gloo --version $(VERSION)
	# hugo --config docs.toml
	# mv site site-versioned

## Old way of deploying site. TODO: remove when old hosting setup is disabled
# .PHONY: deploy-site
# deploy-site: site
# 	firebase deploy --only hosting:gloo-docs

.PHONY: serve-site
serve-site: site-test
	hugo --config docs.toml --themesDir themes server -D

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
	docker build \
		--build-arg VERSION=latest \
		--build-arg PRODUCT_SCOPE=$(PRODUCT_SCOPE) \
		--build-arg FROM_DIR=./site-latest \
		-t $(GCR_REPO_PREFIX)/$(IMAGE_LEAF_NAME):$(SERVE_AS_LATEST_TAG) .
	docker push $(GCR_REPO_PREFIX)/$(IMAGE_LEAF_NAME):$(SERVE_AS_LATEST_TAG)
