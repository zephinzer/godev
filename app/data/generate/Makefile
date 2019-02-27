##
## Makefile constants - extract to a separate file if needed
## ---------------------------------------------------------
## specifies the name of your application binary
BIN_NAME=app
## specifies the relative path to a directory where the binary should be placed in
BIN_PATH=bin
## specifies the registry to push to
DOCKER_REGISTRY_HOSTNAME=docker.io
## specifies docker.io/THIS/image:tag
DOCKER_IMAGE_NAMESPACE=godev
## specifies docker.io/namespace/THIS:tag - align with $(BIN_NAME) for less confusion
DOCKER_IMAGE_NAME=example
## specifies the absolute path to the directory containing the .git directory
GIT_ROOT=$(CURDIR)/../../..
## enable following line to draw variables from a file named Makefile.properties
# include Makefile.properties

## starts the application for development with live-reload
start:
	@godev --watch
## installs the dependencies using go modules
deps:
	@go mod vendor
## runs the tests with live-reload
test:
	@godev --test --watch
## runs the tests just once with coverage
test.once:
	@godev --test
## compiles binaries for all systems
compile:
	@$(MAKE) compile.linux
	@$(MAKE) compile.macos
	@$(MAKE) compile.windows
## compiles binaries for linux
compile.linux:
	@$(MAKE) GOARCH=amd64 GOOS=linux .compile
## compiles binaries for macos
compile.macos:
	@$(MAKE) GOARCH=amd64 GOOS=darwin .compile
## compiles binaries for windows
compile.windows:
	@$(MAKE) GOARCH=386 GOOS=windows BIN_EXT=.exe .compile
## compilation driver
.compile:
	@CGO_EMABLED=0 GO111MODULE=on \
		go build -a -ldflags "-extldflags -static" -o $(CURDIR)/$(BIN_PATH)/$(BIN_NAME)-${GOOS}-${GOARCH}${BIN_EXT}
	@chmod +x $(CURDIR)/$(BIN_PATH)/$(BIN_NAME)-${GOOS}-${GOARCH}${BIN_EXT}
	@sha256sum $(CURDIR)/$(BIN_PATH)/$(BIN_NAME)-${GOOS}-${GOARCH}${BIN_EXT} | cut -d " " -f 1 > $(CURDIR)/$(BIN_PATH)/$(BIN_NAME)-${GOOS}-${GOARCH}${BIN_EXT}.sha256
## dockerisation for production
docker:
	@$(MAKE) .docker STAGE="production"
## dockerisation for development
docker.dev:
	@$(MAKE) .docker STAGE="development"
## dockerisation driver
.docker:
	@$(MAKE) log.info MSG="creating image $(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest"
	@docker build \
		--target ${STAGE} \
		--build-arg BIN_NAME=$(BIN_NAME) \
		--build-arg BIN_PATH=$(BIN_PATH) \
		--target=production \
		-t $(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest \
		.
docker.prepare: docker
	@$(MAKE) log.info MSG="tagging image $(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest"
	@docker tag \
		$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest \
		$(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest
	@$(MAKE) log.info MSG="tagging image $(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*')"
	@docker tag \
		$(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest \
		$(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*')
	@$(MAKE) log.info MSG="tagging image $(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*')-$$(git rev-list -1 HEAD)"
	@docker tag \
		$(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*') \
		$(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*')-$$(git rev-list -1 HEAD)
publish.dockerhub: docker.prepare
	@$(MAKE) log.info MSG="pushing image $(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest"
	@docker push $(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest
	@$(MAKE) log.info MSG="pushing image $(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*')"
	@docker push $(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*')
	@$(MAKE) log.info MSG="pushing image $(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*')-$$(git rev-list -1 HEAD)"
	@docker push $(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*')-$$(git rev-list -1 HEAD)
version.get:
	@docker run \
		-v "$(GIT_ROOT):/app" \
		zephinzer/vtscripts:latest \
		get-latest -q
version.next:
	@docker run \
		-v "$(GIT_ROOT):/app" \
		zephinzer/vtscripts:latest \
		get-next -q
version.bump:
	@docker run \
		-v "$(GIT_ROOT):/app" \
		zephinzer/vtscripts:latest \
		iterate ${VERSION} -i -q
log.debug:
	-@printf -- "\033[36m\033[1m_ [DEBUG] ${MSG}\033[0m\n"
log.info:
	-@printf -- "\033[32m\033[1m>  [INFO] ${MSG}\033[0m\n"
log.warn:
	-@printf -- "\033[33m\033[1m?  [WARN] ${MSG}\033[0m\n"
log.error:
	-@printf -- "\033[31m\033[1m! [ERROR] ${MSG}\033[0m\n"
