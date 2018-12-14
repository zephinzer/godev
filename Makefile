# global variables - change these to your own Docker Hub configuration if 
# you are forking this repository and building your own image out of this
DOCKER_NAMESPACE=zephinzer
DOCKER_IMAGE_NAME=golang-dev

# builds the image for use locally
build:
	@docker build -t $(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest .

# publishes the image to Docker Hub via 3 tags:
# 1. $SEMVER_VERSION tag (if you're fussy about the methodology of this image)
# 2. $GO_VERSION tag (if you're fussy about the golang version but don't care about all else)
# 3. `latest` tag (for development use)
publish: build
	@$(MAKE) publish.app
	@$(MAKE) publish.go
	@docker push $(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest

# publishes the image with the $SEMVER_VERSION tag
publish.app: tag.app
	@printf -- "$$(docker run \
		-v "$(CURDIR):/app" \
		$(DOCKER_NAMESPACE)/vtscripts:latest \
		get-latest -q)" > $(CURDIR)/.app.version
	@docker push $(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$(cat $(CURDIR)/.app.version)
	@rm -rf $(CURDIR)/.app.version

# publishes the image with the $GO_VERSION tag
publish.go: tag.go
	@printf -- "$$(docker run \
		--entrypoint=go \
		$(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest \
		version)" | cut -f 3 -d ' ' | sed -e 's|go||g' > $(CURDIR)/.go.version
	@docker push $(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$(cat $(CURDIR)/.go.version)
	@rm -rf $(CURDIR)/.go.version

# tags the image with the $SEMVER_VERSION tag
tag.app:
	@printf -- "$$(docker run \
		-v "$(CURDIR):/app" \
		$(DOCKER_NAMESPACE)/vtscripts:latest \
		get-latest -q)" > $(CURDIR)/.app.version
	@docker tag $(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest \
		$(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$(cat $(CURDIR)/.app.version)
	@rm -rf $(CURDIR)/.app.version

# tags the image with the $GO_VERSION tag
tag.go:
	@printf -- "$$(docker run \
		--entrypoint=go \
		$(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest \
		version)" | cut -f 3 -d ' ' | sed -e 's|go||g' > $(CURDIR)/.go.version
	@docker tag $(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest \
		$(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$(cat $(CURDIR)/.go.version)
	@rm -rf $(CURDIR)/.go.version

# retrieves the latest tagged version of this repository
version.get:
	@docker run \
		-v "$(CURDIR):/app" \
		zephinzer/vtscripts:latest \
		get-latest -q

# bumps the version of this repository
# to bump the patch version, run this without arguments
# to bump the minor version, run this with VERSION=minor argument
# to bump the major version, run this with VERSION=major argument
version.bump:
	@docker run \
		-v "$(CURDIR):/app" \
		zephinzer/vtscripts:latest \
		iterate ${VERSION} -i

# runs some sanity checks, see ./test for more information
test: build
	cd $(CURDIR)/test && $(MAKE) test
