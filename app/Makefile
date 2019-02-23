GOLANG_DEV_VERSION=latest

# initialises this directory - use once only
init:
	@$(MAKE) _dev ARG="init"
# builds the application - outputs an `app` binary
build: 
	@$(MAKE) _dev ARG="build" ARGS=
# runs tests in watch-mode
test: build
	@$(MAKE) _dev ARG="test"
# runs tests once
test.once: build
	@$(MAKE) _dev ARG="test -coverprofile c.out"
# runs the application on the host network
start: build
	@$(MAKE) _dev ARG="start" ARGS=${ARGS}
# creates a shell in a fresh container generated from the image, usable for development on non-linux machines
shell:
	$(MAKE) _dev ARG="shell"
# retrieves the latest version we are at
version.get:
	@docker run -v "$(CURDIR):/app" zephinzer/vtscripts:latest get-latest -q
# bumps the version by 1: specify VERSION as "patch", "minor", or "major", to be specific about things
version.bump: 
	@docker run -v "$(CURDIR):/app" zephinzer/vtscripts:latest iterate ${VERSION} -i
# base command to run other scripts (do not use alone)
_dev:
	@docker run \
    -it \
    --network host \
    -u $$(id -u) \
    -v "$(CURDIR)/.cache/pkg:/go/pkg" \
    -v "$(CURDIR):/go/src/app" \
    zephinzer/golang-dev:$(GOLANG_DEV_VERSION) ${ARG} ${ARGS}
