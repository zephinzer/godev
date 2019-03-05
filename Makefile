# relative path to the binary directory - remove any directory separators
# from the start and end of the variable before use
BINARY_PATH=bin
# name of the binary filename
BINARY_FILENAME=godev
# THIS/zephinzer/godev:latest
DOCKER_IMAGE_REGISTRY=docker.io
# docker.io/THIS/godev:latest
DOCKER_IMAGE_NAMESPACE=zephinzer
# docker.io/zephinzer/THIS:latest
DOCKER_IMAGE_NAME=godev

##
# - VERSION
# - COMMIT
compile:
	go build \
		-a \
		-o $(CURDIR)/$(BINARY_PATH)/$(BINARY_FILENAME) \
		-ldflags " \
			-extldflags -static \
			-X main.Version=${VERSION} \
			-X main.Commit=${COMMIT} \
		"
docker:
	@mkdir -p $(CURDIR)/.docker
	@printf -- 'if this .docker directory is here, it means "make docker" was terminated unexpectedly' > $(CURDIR)/.docker/README
	@printf -- "$$(git rev-list -1 HEAD | head -c 7)" > $(CURDIR)/.docker/.Commit
	@printf -- "$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*')" > $(CURDIR)/.docker/.Version
	@docker build \
		--build-arg VERSION=$$(cat $(CURDIR)/.docker/.Version) \
		--build-arg COMMIT=$$(cat $(CURDIR)/.docker/.Commit) \
		-t $(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest \
		.
	@docker tag $(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest $(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$(cat $(CURDIR)/.docker/.Version)
	@docker run $(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest go version > $(CURDIR)/.docker/.GoVersion
	@docker tag $(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest $(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$(cat $(CURDIR)/.docker/.Version)
	@rm -rf $(CURDIR)/.docker
release.docker: docker
	@docker tag 
## generates the contributors file
contributors:
	@echo "# generate with 'make contributors'\n#" > $(CURDIR)/CONTRIBUTORS
	@echo "# last generated on $$(date -u)\n" >> $(CURDIR)/CONTRIBUTORS
	@git shortlog -se | sed -e 's|@|-at-|g' -e 's|\.|-dot-|g' | cut -f 2- >> $(CURDIR)/CONTRIBUTORS
## retrieves the latest version we are at
version.get:
	@docker run -v "$(CURDIR):/app" zephinzer/vtscripts:latest get-latest -q
## bumps the version by 1: specify VERSION as "patch", "minor", or "major", to be specific about things
version.bump: 
	@docker run -v "$(CURDIR):/app" zephinzer/vtscripts:latest iterate ${VERSION} -i
## blue logs
log.debug:
	-@printf -- "\033[36m\033[1m_ [DEBUG] ${MSG}\033[0m\n"
## green logs
log.info:
	-@printf -- "\033[32m\033[1m>  [INFO] ${MSG}\033[0m\n"
## yellow logs
log.warn:
	-@printf -- "\033[33m\033[1m?  [WARN] ${MSG}\033[0m\n"
## red logs (die mf)
log.error:
	-@printf -- "\033[31m\033[1m! [ERROR] ${MSG}\033[0m\n"
