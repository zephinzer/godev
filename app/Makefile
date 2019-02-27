GOLANG_DEV_VERSION=latest

compile: generate
	@$(MAKE) log.info MSG="generating static binary..."
	@CGO_ENABLED=0 GO111MODULE=on \
		go build \
			-a \
			-o $(CURDIR)/bin/godev \
			-ldflags " \
				-extldflags -static \
				-X main.Version=$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*') \
				-X main.Commit=$$(git rev-list -1 HEAD) \
			"
	@$(MAKE) log.info MSG="generated binary at $(CURDIR)/bin/godev"
generate:
	@$(MAKE) log.info MSG="generating static data file data.go (see ./data/generate.go)..."
	@go generate
	@$(MAKE) log.info MSG="generated data.go..."
start: generate
	@$(MAKE) log.info MSG="running application for development with live-reload..."
	@$(MAKE) _dev ARG="start" ARGS="--test --ignore bin,data,vendor,.cache ${ARGS}"
start.prd: generate
	@$(MAKE) log.info MSG="running application for development in production with live-reload..."
	@$(MAKE) _dev ARG="start" ARGS="${ARGS}"
run: install.deps generate
	@$(MAKE) log.info MSG="running application in development..."
	@go run $$(ls -a | grep .go | grep -v "test" | tr -s '\n' ' ') ${ARGS}
run.prd: install.deps compile
	@$(MAKE) log.info MSG="running application in production..."
	@bin/godev ${ARGS}
install.deps:
	@$(MAKE) log.info MSG="installing dependencies with go modules..."
	@GO111MODULE=on go mod vendor
test: install.deps generate
	@$(MAKE) log.info MSG="running tests with live-reload"
	@$(MAKE) _dev ARG="test"
test.once: install.deps generate
	@$(MAKE) log.info MSG="running tests once with coverage output"
	@$(MAKE) _dev ARG="test -coverprofile c.out"
shell:
	@$(MAKE) log.info MSG="creating a shell in the docker development environment..."
	$(MAKE) _dev ARG="shell"
## retrieves the latest version we are at
version.get:
	@docker run -v "$(CURDIR)/..:/app" zephinzer/vtscripts:latest get-latest -q
## bumps the version by 1: specify VERSION as "patch", "minor", or "major", to be specific about things
version.bump: 
	@docker run -v "$(CURDIR)/..:/app" zephinzer/vtscripts:latest iterate ${VERSION} -i
## driver recipe to run other scripts (do not use alone)
_dev:
	@docker run \
    -it \
    --network host \
    -u $$(id -u) \
    -v "$(CURDIR)/.cache/pkg:/go/pkg" \
    -v "$(CURDIR):/go/src/app" \
    zephinzer/golang-dev:$(GOLANG_DEV_VERSION) ${ARG} ${ARGS}
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
