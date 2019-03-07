include Makefile.properties

## compiles the godev binary for all platforms using Docker
## - call this when running on host
godev:
	@$(MAKE) create.version.data FOR_OP=godev
	@docker build \
		--target=build \
		--build-arg VERSION=$$(cat ./.godev/.Version) \
		--build-arg COMMIT=$$(cat ./.godev/.Commit) \
		-t $(DOCKER_IMAGE_NAME)_$$(cat ./.godev/.Commit):latest \
		.
	@$(MAKE) log.debug MSG="terminating any existing instances of '$(DOCKER_IMAGE_NAME)_$$(cat ./.godev/.Commit)'..."
	-@docker stop $(DOCKER_IMAGE_NAME)_$$(cat ./.godev/.Commit)
	-@docker rm $(DOCKER_IMAGE_NAME)_$$(cat ./.godev/.Commit)
	@$(MAKE) log.debug MSG="starting an instance of '$(DOCKER_IMAGE_NAME)_$$(cat ./.godev/.Commit)'..."
	@docker run -d \
		--entrypoint=sleep \
		--name=$(DOCKER_IMAGE_NAME)_$$(cat ./.godev/.Commit) \
		$(DOCKER_IMAGE_NAME)_$$(cat ./.godev/.Commit):latest \
		1000
	@$(MAKE) log.debug MSG="copying out binaries from '$(DOCKER_IMAGE_NAME)_$$(cat ./.godev/.Commit)' to '$(CURDIR)/bin'..."
	@docker exec $(DOCKER_IMAGE_NAME)_$$(cat ./.godev/.Commit) ls -1 /go/build/bin | xargs -I @ docker cp $(DOCKER_IMAGE_NAME)_$$(cat ./.godev/.Commit):/go/build/bin/@ $(CURDIR)/bin/@
	@$(MAKE) log.debug MSG="terminating '$(DOCKER_IMAGE_NAME)_$$(cat ./.godev/.Commit)'..."
	@docker stop $(DOCKER_IMAGE_NAME)_$$(cat ./.godev/.Commit)
	@docker rm $(DOCKER_IMAGE_NAME)_$$(cat ./.godev/.Commit)
	@rm -rf ./.godev
## compiles the godev binary for all platforms on the host
## - call this directly when running for docker build
compile:
	@$(MAKE) log.debug MSG="compiling godev..."
	@$(MAKE) VERSION=${VERSION} COMMIT=${COMMIT} compile.linux
	@$(MAKE) VERSION=${VERSION} COMMIT=${COMMIT} compile.macos
	@$(MAKE) VERSION=${VERSION} COMMIT=${COMMIT} compile.windows
## compiles the binary for linux
compile.linux:
	@$(MAKE) log.debug MSG="compiling godev for linux..."
	@$(MAKE) VERSION=${VERSION} COMMIT=${COMMIT} GOOS=linux GOARCH=amd64 .compile
## compiles the binary for macos
compile.macos:
	@$(MAKE) log.debug MSG="compiling godev for macos..."
	@$(MAKE) VERSION=${VERSION} COMMIT=${COMMIT} GOOS=darwin GOARCH=amd64 .compile
## compiles the binary for windows
compile.windows:
	@$(MAKE) log.debug MSG="compiling godev for windows..."
	@$(MAKE) VERSION=${VERSION} COMMIT=${COMMIT} GOOS=windows GOARCH=386 BINARY_EXT=.exe .compile
## generic compilation recipe for ensuring consistency of above recipes
.compile: deps
	@$(MAKE) log.debug MSG="compiling godev at $(CURDIR)/bin/godev-${VERSION}-${GOOS}-${GOARCH}${BINARY_EXT} - version: '${VERSION}' commit: '${COMMIT}'..."
	@CGO_ENABLED=0 \
		GO111MODULES=on \
		GOOS=${GOOS} \
		GOARCH=${GOARCH} \
		go build \
		-a \
		-o $(CURDIR)/bin/godev-${VERSION}-${GOOS}-${GOARCH}${BINARY_EXT} \
		-ldflags " \
			-extldflags -static \
		"
	@sha256sum $(CURDIR)/bin/godev-${VERSION}-${GOOS}-${GOARCH}${BINARY_EXT} | cut -d ' ' -f 1 > $(CURDIR)/bin/godev-${VERSION}-${GOOS}-${GOARCH}${BINARY_EXT}.sha256
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
test.ci: deps
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
	@$(MAKE) log.debug MSG="building docker image '$(DOCKER_IMAGE_NAME):latest'..."
	# build docker image
	@docker build \
		--build-arg VERSION=$$(cat ./.docker/.Version) \
		--build-arg COMMIT=$$(cat ./.docker/.Commit) \
		-t $(DOCKER_IMAGE_NAME):latest \
		.
	# retrieve Go version
	@docker run $(DOCKER_IMAGE_NAME):latest go version | sed 's|go||g' | cut -f 3 -d ' ' \
		> ./.docker/.GoVersion
	@$(MAKE) _docker.tag TAG="latest"
	@$(MAKE) _docker.tag TAG="$$(cat ./.docker/.Version)"
	@$(MAKE) _docker.tag TAG="$$(cat ./.docker/.Commit)"
	@$(MAKE) _docker.tag TAG="$$(cat ./.docker/.Version)-$$(cat ./.docker/.Commit)"
	@$(MAKE) _docker.tag TAG="go-$$(cat ./.docker/.GoVersion)"
	@rm -rf ./.docker
_docker.tag:
	@$(MAKE) log.debug MSG="tagging docker image '$$(cat ./.docker/.DockerImage):${TAG}'..."
	@docker tag $(DOCKER_IMAGE_NAME):latest $$(cat ./.docker/.DockerImage):${TAG}

## releases the docker image
release.docker: docker
	@$(MAKE) create.version.data FOR_OP=release.docker
	@docker run $$(cat ./.release.docker/.DockerImage):latest go version | sed 's|go||g' | cut -f 3 -d ' ' \
		> ./.release.docker/.GoVersion
	@$(MAKE) _release.docker.push TAG="latest"
	@$(MAKE) _release.docker.push TAG="$$(cat ./.release.docker/.Version)"
	@$(MAKE) _release.docker.push TAG="$$(cat ./.release.docker/.Commit)"
	@$(MAKE) _release.docker.push TAG="$$(cat ./.release.docker/.Version)-$$(cat ./.release.docker/.Commit)"
	@$(MAKE) _release.docker.push TAG="go-$$(cat ./.release.docker/.GoVersion)"
	@rm -rf ./.release.docker
_release.docker.push:
	@$(MAKE) logs.debug MSG="pushing '$$(cat ./.release.docker/.DockerImage):${TAG}'..."
	@docker push $$(cat ./.release.docker/.DockerImage):${TAG}

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
	@mkdir -p ./.${FOR_OP}
	@printf -- 'if this .${FOR_OP} directory is here, it means "make ${FOR_OP}" was terminated unexpectedly' > ./.${FOR_OP}/README
	@$(MAKE) log.debug MSG="generating git commit..."
	@printf -- "$$(git rev-list -1 HEAD | head -c 7)" \
		> ./.${FOR_OP}/.Commit
	@$(MAKE) log.debug MSG="generating git tag..."
	@printf -- "$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*')" \
		> ./.${FOR_OP}/.Version
	@$(MAKE) log.debug MSG="generating docker image name..."
	@printf -- '$(DOCKER_IMAGE_REGISTRY)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME)' \
		> ./.${FOR_OP}/.DockerImage

## creates a set of keys you can use for deploy keys
ssh.keys: # PREFIX= - defaults to nothing if not specified
	@ssh-keygen -t rsa -b 8192 -f ./bin/${PREFIX}_id_rsa -q -N ''
	@cat ./bin/${PREFIX}_id_rsa | base64 -w 0 > ./bin/${PREFIX}_id_rsa_b64

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
