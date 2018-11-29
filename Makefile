DOCKER_NAMESPACE=zephinzer
DOCKER_IMAGE_NAME=golang

build:
	@docker build -t $(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest .

publish: build
	@$(MAKE) publish.app
	@$(MAKE) publish.go

publish.app: tag.app
	@printf -- "$$(docker run \
		-v "$(CURDIR):/app" \
		$(DOCKER_NAMESPACE)/vtscripts:latest \
		get-latest -q)" > $(CURDIR)/.app.version
	@docker push $(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$(cat $(CURDIR)/.app.version)
	@rm -rf $(CURDIR)/.app.version

publish.go: tag.go
	@printf -- "$$(docker run \
		--entrypoint=go \
		$(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest \
		version)" | cut -f 3 -d ' ' | sed -e 's|go||g' > $(CURDIR)/.go.version
	@docker push $(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$(cat $(CURDIR)/.go.version)
	@rm -rf $(CURDIR)/.go.version

tag.app:
	@printf -- "$$(docker run \
		-v "$(CURDIR):/app" \
		$(DOCKER_NAMESPACE)/vtscripts:latest \
		get-latest -q)" > $(CURDIR)/.app.version
	@docker tag $(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest \
		$(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$(cat $(CURDIR)/.app.version)
	@rm -rf $(CURDIR)/.app.version

tag.go:
	@printf -- "$$(docker run \
		--entrypoint=go \
		$(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest \
		version)" | cut -f 3 -d ' ' | sed -e 's|go||g' > $(CURDIR)/.go.version
	@docker tag $(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest \
		$(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$(cat $(CURDIR)/.go.version)
	@rm -rf $(CURDIR)/.go.version

version.get:
	@docker run \
		-v "$(CURDIR):/app" \
		zephinzer/vtscripts:latest \
		get-latest -q

version.bump:
	@docker run \
		-v "$(CURDIR):/app" \
		zephinzer/vtscripts:latest \
		iterate ${VERSION} -i

test: build
	cd ./test && $(MAKE) test
