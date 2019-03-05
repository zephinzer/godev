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
	@$(MAKE) log.debug MSG="compiling godev at $(CURDIR)/$(BINARY_PATH)/$(BINARY_FILENAME) - version: '${VERSION}' commit: '${COMMIT}'..."
	@go build \
		-a \
		-o $(CURDIR)/$(BINARY_PATH)/$(BINARY_FILENAME) \
		-ldflags " \
			-extldflags -static \
			-X main.Version=${VERSION} \
			-X main.Commit=${COMMIT} \
		"
	@$(MAKE) log.info MSG="compiled godev at $(CURDIR)/$(BINARY_PATH)/$(BINARY_FILENAME) - version: '${VERSION}' commit: '${COMMIT}'"
start: compile
	@$(MAKE) log.debug MSG="running godev for development..."
	@$(CURDIR)/$(BINARY_PATH)/$(BINARY_FILENAME) -vv --watch $(CURDIR) --dir $(CURDIR)/dev
deps:
	@$(MAKE) log.debug MSG="installing the dependencies..."
	@go mod vendor
	@$(MAKE) log.info MSG="dependency installation successful."
generate:
	@$(MAKE) log.debug MSG="generating a new ~/data.go..."
	@go generate
	@$(MAKE) log.info MSG="~/data.go generation successful."
test: compile
	@$(MAKE) log.debug MSG="running tests in watch mode for godev..."
	@$(CURDIR)/$(BINARY_PATH)/$(BINARY_FILENAME) --test
binary:

docker:
	@$(MAKE) create.version.data FOR_OP=docker
	@printf -- '$(DOCKER_IMAGE_REGISTRY)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME)' \
		> $(CURDIR)/.docker/.DockerImage
	@$(MAKE) log.debug MSG="building docker image '$$(cat $(CURDIR)/.docker/.DockerImage):latest'..."
	# build docker image
	@docker build \
		--build-arg VERSION=$$(cat $(CURDIR)/.docker/.Version) \
		--build-arg COMMIT=$$(cat $(CURDIR)/.docker/.Commit) \
		-t $$(cat $(CURDIR)/.docker/.DockerImage):latest \
		.
	# tag version
	@$(MAKE) log.debug MSG="tagging docker image '$$(cat $(CURDIR)/.docker/.DockerImage):$$(cat $(CURDIR)/.docker/.Version)'..."
	@docker tag $$(cat $(CURDIR)/.docker/.DockerImage):latest \
		$$(cat $(CURDIR)/.docker/.DockerImage):$$(cat $(CURDIR)/.docker/.Version)
	# tag commit
	@$(MAKE) log.debug MSG="tagging docker image '$$(cat $(CURDIR)/.docker/.DockerImage):$$(cat $(CURDIR)/.docker/.Commit)'..."
	@docker tag $$(cat $(CURDIR)/.docker/.DockerImage):latest \
		$$(cat $(CURDIR)/.docker/.DockerImage):$$(cat $(CURDIR)/.docker/.Commit)
	# tag version-commit
	@$(MAKE) log.debug MSG="tagging docker image '$$(cat $(CURDIR)/.docker/.DockerImage):$$(cat $(CURDIR)/.docker/.Version)-$$(cat $(CURDIR)/.docker/.Commit)'..."
	@docker tag $$(cat $(CURDIR)/.docker/.DockerImage):latest \
		$$(cat $(CURDIR)/.docker/.DockerImage):$$(cat $(CURDIR)/.docker/.Version)-$$(cat $(CURDIR)/.docker/.Commit)
	@docker run $$(cat $(CURDIR)/.docker/.DockerImage):latest go version | sed 's|go||g' | cut -f 3 -d ' ' \
		> $(CURDIR)/.docker/.GoVersion
	# tag go version
	@$(MAKE) log.debug MSG="tagging docker image '$$(cat $(CURDIR)/.docker/.DockerImage):go-$$(cat $(CURDIR)/.docker/.GoVersion)'..."
	@docker tag $$(cat $(CURDIR)/.docker/.DockerImage):latest \
		$$(cat $(CURDIR)/.docker/.DockerImage):go-$$(cat $(CURDIR)/.docker/.GoVersion)
	@rm -rf $(CURDIR)/.docker
release.docker: docker
	@$(MAKE) create.version.data FOR_OP=release.docker
	@printf -- '$(DOCKER_IMAGE_REGISTRY)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME)' \
		> $(CURDIR)/.release.docker/.DockerImage
	@docker run $$(cat $(CURDIR)/.release.docker/.DockerImage):latest go version | sed 's|go||g' | cut -f 3 -d ' ' \
		> $(CURDIR)/.release.docker/.GoVersion
	# push everything we built in 'make docker'
	@docker push$$(cat $(CURDIR)/.docker/.DockerImage):latest
	@docker push $$(cat $(CURDIR)/.docker/.DockerImage):$$(cat $(CURDIR)/.docker/.Verson)
	@docker push $$(cat $(CURDIR)/.docker/.DockerImage):$$(cat $(CURDIR)/.docker/.Commit)
	@docker push $$(cat $(CURDIR)/.docker/.DockerImage):$$(cat $(CURDIR)/.docker/.Version)-$$(cat $(CURDIR)/.docker/.Commit)
	@docker push $$(cat $(CURDIR)/.docker/.DockerImage):go-$$(cat $(CURDIR)/.docker/.GoVersion)
## generates the contributors file
contributors:
	@echo "# generate with 'make contributors'\n#" > $(CURDIR)/CONTRIBUTORS
	@echo "# last generated on $$(date -u)\n" >> $(CURDIR)/CONTRIBUTORS
	@git shortlog -se | sed -e 's|@|-at-|g' -e 's|\.|-dot-|g' | cut -f 2- >> $(CURDIR)/CONTRIBUTORS
create.version.data:
	@$(MAKE) log.debug MSG="provisioning '$(CURDIR)/.${FOR_OP}' directory..."
	@mkdir -p $(CURDIR)/.${FOR_OP}
	@printf -- 'if this .${FOR_OP} directory is here, it means "make ${FOR_OP}" was terminated unexpectedly' > $(CURDIR)/.${FOR_OP}/README
	@$(MAKE) log.debug MSG="generating git commit..."
	@printf -- "$$(git rev-list -1 HEAD | head -c 7)" > $(CURDIR)/.${FOR_OP}/.Commit
	@$(MAKE) log.debug MSG="generating git tag..."
	@printf -- "$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*')" > $(CURDIR)/.${FOR_OP}/.Version
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
