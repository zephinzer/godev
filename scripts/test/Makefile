test:
	@$(MAKE) test.base
	@$(MAKE) test.run
	@$(MAKE) test.files
	@$(MAKE) test.build
	@$(MAKE) test.usage.build
	@$(MAKE) test.usage.run

test.base:
	@$(MAKE) log.debug MSG='building the base image...';
	@docker build -f ../Dockerfile -t zephinzer/golang-dev:test ..
	@$(MAKE) log.info MSG='base image built successfully.';

test.build:
	@$(MAKE) log.debug MSG='testing the build...'
	-@rm -rf $(CURDIR)/app
	@docker run \
		-v "$(CURDIR):/go/src/app" \
		--entrypoint=build \
		zephinzer/golang-dev:test
	@if [ ! -e "$(CURDIR)/app" ]; then exit 1; else rm -rf $(CURDIR)/app; fi
	@$(MAKE) log.info MSG='application successfully built.';

test.files:
	@$(MAKE) log.debug MSG='making sure /scripts exists...';
	@docker run \
		-v "$(CURDIR):/go/src/app" \
		--entrypoint=ls \
		zephinzer/golang-dev:test \
		-1 /scripts
	@$(MAKE) log.info MSG='/scripts exists.';

test.run:
	@$(MAKE) log.debug MSG='running the base image...';
	@docker run \
		-v "$(CURDIR):/go/src/app" \
		zephinzer/golang-dev:test
	@$(MAKE) log.info MSG='base image ran successfully.';

test.usage.build:
	@$(MAKE) log.debug MSG='using base image to build sub-image with scratch...'
	@docker build -t zephinzer/golang-dev:testusage .
	@$(MAKE) log.info MSG='scratch image built successfully.'

test.usage.run:
	@$(MAKE) log.debug MSG='running image based on scratch...'
	-@docker stop zephinzer-golang-testusage
	-@docker rm zephinzer-golang-testusage
	@docker run -d --name zephinzer-golang-testusage zephinzer/golang-dev:testusage
	@docker cp zephinzer-golang-testusage:/.output $(CURDIR)/.output
	-@docker stop zephinzer-golang-testusage
	-@docker rm zephinzer-golang-testusage
	@cat $(CURDIR)/.output | grep -e "It works!"
	@rm -rf $(CURDIR)/.output
	$(MAKE) log.info MSG='scratch image ran successfully.'

log.debug:
	-@printf -- "\033[36m\033[1m_ [DEBUG] ${MSG}\033[0m\n"
log.info:
	-@printf -- "\033[32m\033[1m>  [INFO] ${MSG}\033[0m\n"
log.warn:
	-@printf -- "\033[33m\033[1m?  [WARN] ${MSG}\033[0m\n"
log.error:
	-@printf -- "\033[31m\033[1m! [ERROR] ${MSG}\033[0m\n"

eg.test:
	docker run -it -u $$(id -u) -v "$(CURDIR):/go/src/app" zephinzer/golang-dev:latest test
eg.start:
	docker run -it -u $$(id -u) -v "$(CURDIR):/go/src/app" zephinzer/golang-dev:latest start
eg.build:
	docker run -it -u $$(id -u) -v "$(CURDIR):/go/src/app" zephinzer/golang-dev:latest build
