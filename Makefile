# relative path to the binary directory - remove any directory separators
# from the start and end of the variable before use
BINARY_PATH=bin
# name of the binary filename
BINARY_FILENAME=godev
GOLANG_DEV_VERSION=latest

compile:
	@$(MAKE) compile.linux
	@$(MAKE) compile.macos
	@$(MAKE) compile.windows
compile.linux:
	@$(MAKE) GOARCH=amd64 GOOS=linux .compile
compile.macos:
	@$(MAKE) GOARCH=amd64 GOOS=darwin .compile
compile.windows:
	@$(MAKE) GOARCH=386 GOOS=windows BINARY_EXT=.exe .compile
## run this to generate the binary
# - BINARY_EXT
# - GOARCH
# - GOOS
.compile: deps generate
	@$(MAKE) log.info MSG="generating static binary..."
	@mkdir -p $(CURDIR)/.compile
	@printf -- 'if this file is here, a compilation is in process/has errored out' > $(CURDIR)/.compile/README
	@printf -- "$$(git rev-list -1 HEAD | head -c 7)" > $(CURDIR)/.compile/.Commit
	@printf -- "$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*')" > $(CURDIR)/.compile/.Version
	@printf -- "$(CURDIR)/$(BINARY_PATH)/$(BINARY_FILENAME)${BINARY_EXT}" > $(CURDIR)/.compile/.bin
	@printf -- "$(CURDIR)/$(BINARY_PATH)/$(BINARY_FILENAME)-${GOOS}-${GOARCH}${BINARY_EXT}" > $(CURDIR)/.compile/.binarch
	@printf -- "$(CURDIR)/$(BINARY_PATH)/$(BINARY_FILENAME)-$$(cat $(CURDIR)/.compile/.Version)-${GOOS}-${GOARCH}${BINARY_EXT}" > $(CURDIR)/.compile/.binver
	@printf -- "$(CURDIR)/$(BINARY_PATH)/$(BINARY_FILENAME)-$$(cat $(CURDIR)/.compile/.Version)-$$(git rev-list -1 HEAD | head -c 7)-${GOOS}-${GOARCH}${BINARY_EXT}" > $(CURDIR)/.compile/.binext
	@CGO_ENABLED=0 \
		GO111MODULE=on \
		GOARCH=${GOARCH} \
		GOOS=${GOOS} \
		go build \
			-a \
			-o $$(cat $(CURDIR)/.compile/.binext) \
			-ldflags " \
				-extldflags -static \
				-X main.Version=$$(cat $(CURDIR)/.compile/.Version) \
				-X main.Commit=$$(cat $(CURDIR)/.compile/.Commit) \
			"
	@chmod +x $$(cat $(CURDIR)/.compile/.binext)
	@ln -s $$(cat $(CURDIR)/.compile/.binext) $$(cat $(CURDIR)/.compile/.bin)
	@ln -s $$(cat $(CURDIR)/.compile/.binext) $$(cat $(CURDIR)/.compile/.binarch)
	@ln -s $$(cat $(CURDIR)/.compile/.binext) $$(cat $(CURDIR)/.compile/.binver)
	@cp 
	@rm -rf $(CURDIR)/.compile
	@$(MAKE) log.info MSG="generated binary at $(CURDIR)/$(BINARY_PATH)/$(BINARY_FILENAME)${BINARY_EXT}"
start:
	@$(MAKE) start.dev
start.dev:
	@$(MAKE) log.info MSG="running application in development (watching application at $(CURDIR)/dev)..."
	@$(MAKE) .start ARGS="-vvv --dir $(CURDIR)/dev"
start.prd:
	@$(MAKE) log.info MSG="running application in production (watching application at $(CURDIR)/dev)..."
	@$(MAKE) .start ARGS="--dir $(CURDIR)/dev"
start.test:
	@$(MAKE) log.info MSG="running tests..."
	@$(MAKE) .start ARGS="--test -vvv --ignore .cache,.vscode,bin,data,docs,scripts,vendor"
.start: deps generate
	@$(MAKE) log.info MSG="running application..."
	@go run $$(ls -a | grep .go | grep -v "test" | tr -s '\n' ' ') ${ARGS}
generate:
	@$(MAKE) log.info MSG="generating static data file data.go (see ./data/generate.go)..."
	@go generate
	@$(MAKE) log.info MSG="generated data.go..."
deps:
	@$(MAKE) log.info MSG="installing dependencies with go modules..."
	@GO111MODULE=on go mod vendor
docker:
	
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
