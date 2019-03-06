include Makefile.properties

## creates the godev binary for all platforms
binary:
	@$(MAKE) create.version.data FOR_OP=binary
	@$(MAKE) \
		VERSION=$$(cat $(CURDIR)/.binary/.Version) \
		COMMIT=$$(cat $(CURDIR)/.binary/.Commit) \
		compile
	@rm -rf $(CURDIR)/.binary
## compiles the binary for all platforms
## - driver for the `binary` recipe
compile:
	@$(MAKE) VERSION=${VERSION} COMMIT=${COMMIT} compile.linux
	@$(MAKE) VERSION=${VERSION} COMMIT=${COMMIT} compile.macos
	@$(MAKE) VERSION=${VERSION} COMMIT=${COMMIT} compile.windows
## compiles the binary for linux
compile.linux:
	@$(MAKE) VERSION=${VERSION} COMMIT=${COMMIT} GOOS=linux GOARCH=amd64 .compile
## compiles the binary for macos
compile.macos:
	@$(MAKE) VERSION=${VERSION} COMMIT=${COMMIT} GOOS=darwin GOARCH=amd64 .compile
## compiles the binary for windows
compile.windows:
	@$(MAKE) VERSION=${VERSION} COMMIT=${COMMIT} GOOS=windows GOARCH=386 BINARY_EXT=.exe .compile
## generic compilation recipe for ensuring consistency of above recipes
.compile: deps
	@$(MAKE) log.debug MSG="compiling godev at ./bin/godev-${VERSION}-${GOOS}-${GOARCH}${BINARY_EXT} - version: '${VERSION}' commit: '${COMMIT}'..."
	@CGO_ENABLED=0 \
		GO111MODULES=on \
		GOOS=${GOOS} \
		GOARCH=${GOARCH} \
		go build \
		-a \
		-o $(CURDIR)/bin/godev-${VERSION}-${GOOS}-${GOARCH}${BINARY_EXT} \
		-ldflags " \
			-extldflags -static \
			-X main.Version=${VERSION} \
			-X main.Commit=${COMMIT} \
		"
	@$(MAKE) log.info MSG="compiled godev at ./bin/godev-${VERSION}-${GOOS}-${GOARCH}${BINARY_EXT} - version: '${VERSION}' commit: '${COMMIT}'"

## starts the application for development
## - note 1: there is no watch mode, restart this for any changes made
## - note 2: the working directory is at ./dev which runs a server that prints 'hello world' on start
start: compile
	@$(MAKE) log.debug MSG="running godev for development..."
	@$(CURDIR)/bin/godev -vv --watch $(CURDIR) --dir $(CURDIR)/dev

## installs the dependencies
deps:
	@$(MAKE) log.debug MSG="installing the dependencies..."
	@GO111MODULE=on go mod vendor
	@$(MAKE) log.info MSG="dependency installation successful."

## generates the static files
generate:
	@$(MAKE) log.debug MSG="generating a new ~/data.go..."
	@go generate
	@$(MAKE) log.info MSG="~/data.go generation successful."

## runs tests in watch mode
test: compile
	@$(MAKE) log.debug MSG="running tests in watch mode for godev..."
	@$(CURDIR)/bin/godev --test

## runs tests for ci
test.ci: deps compile
	@$(MAKE) log.debug MSG="running tests in single run mode..."
	@go test ./... -coverprofile c.out

## generates the contributors file
contributors:
	@echo "# generate with 'make contributors'\n#" > $(CURDIR)/CONTRIBUTORS
	@echo "# last generated on $$(date -u)\n" >> $(CURDIR)/CONTRIBUTORS
	@git shortlog -se | sed -e 's|@|-at-|g' -e 's|\.|-dot-|g' | cut -f 2- >> $(CURDIR)/CONTRIBUTORS

## builds the docker image
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

## releases the docker image
release.docker: docker
	@$(MAKE) create.version.data FOR_OP=release.docker
	@printf -- '$(DOCKER_IMAGE_REGISTRY)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME)' \
		> $(CURDIR)/.release.docker/.DockerImage
	@docker run $$(cat $(CURDIR)/.release.docker/.DockerImage):latest go version | sed 's|go||g' | cut -f 3 -d ' ' \
		> $(CURDIR)/.release.docker/.GoVersion
	# push latest
	@$(MAKE) logs.debug MSG="pushing '$$(cat $(CURDIR)/.docker/.DockerImage):latest'..."
	@docker push $$(cat $(CURDIR)/.docker/.DockerImage):latest
	# push version
	@$(MAKE) logs.debug MSG="pushing '$$(cat $(CURDIR)/.docker/.DockerImage):$$(cat $(CURDIR)/.docker/.Verson)'..."
	@docker push $$(cat $(CURDIR)/.docker/.DockerImage):$$(cat $(CURDIR)/.docker/.Verson)
	# push commit
	@$(MAKE) logs.debug MSG="pushing '$$(cat $(CURDIR)/.docker/.DockerImage):$$(cat $(CURDIR)/.docker/.Commit)'..."
	@docker push $$(cat $(CURDIR)/.docker/.DockerImage):$$(cat $(CURDIR)/.docker/.Commit)
	# push version-commit
	@$(MAKE) logs.debug MSG="pushing '$$(cat $(CURDIR)/.docker/.DockerImage):$$(cat $(CURDIR)/.docker/.Version)-$$(cat $(CURDIR)/.docker/.Commit)'..."
	@docker push $$(cat $(CURDIR)/.docker/.DockerImage):$$(cat $(CURDIR)/.docker/.Version)-$$(cat $(CURDIR)/.docker/.Commit)
	# push go version
	@$(MAKE) logs.debug MSG="pushing '$$(cat $(CURDIR)/.docker/.DockerImage):go-$$(cat $(CURDIR)/.docker/.GoVersion)'..."
	@docker push $$(cat $(CURDIR)/.docker/.DockerImage):go-$$(cat $(CURDIR)/.docker/.GoVersion)
	@rm -rf $(CURDIR)/.release.docker

## releases tags to github
release.github:
	@if [ "${GITHUB_REPOSITORY_URL}" = "" ]; then exit 1; fi;
	@git remote set-url origin $(GITHUB_REPOSITORY_URL)
	@git checkout --f master
	@git fetch
	@$(MAKE) version.get
	@$(MAKE) version.bump VERSION=${BUMP}
	@$(MAKE) version.get
	@git push --tags

## creates versioning data for use when releasing
create.version.data:
	@$(MAKE) log.debug MSG="provisioning '$(CURDIR)/.${FOR_OP}' directory..."
	@mkdir -p $(CURDIR)/.${FOR_OP}
	@printf -- 'if this .${FOR_OP} directory is here, it means "make ${FOR_OP}" was terminated unexpectedly' > $(CURDIR)/.${FOR_OP}/README
	@$(MAKE) log.debug MSG="generating git commit..."
	@printf -- "$$(git rev-list -1 HEAD | head -c 7)" \
		> $(CURDIR)/.${FOR_OP}/.Commit
	@$(MAKE) log.debug MSG="generating git tag..."
	@printf -- "$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*')" \
		> $(CURDIR)/.${FOR_OP}/.Version

## creates a set of keys you can use for deploy keys
ssh.keys: # PREFIX= - defaults to nothing if not specified
	@ssh-keygen -t rsa -b 8192 -f $(CURDIR)/bin/${PREFIX}_id_rsa -q -N ''
	@cat $(CURDIR)/bin/${PREFIX}_id_rsa | base64 -w 0 > $(CURDIR)/bin/${PREFIX}_id_rsa_b64

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
