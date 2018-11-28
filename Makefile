build:
	$(MAKE) version.get > .version
	# docker build -t zephinzer/golang:latest .

version.get:
	docker run \
		-v "$(CURDIR):/app" \
		zephinzer/vtscripts:latest \
		get-latest -h

version.bump:
	docker run \
		-v "$(CURDIR):/app" \
		zephinzer/vtscripts:latest \
		iterate patch -q

test: build
	cd ./test && $(MAKE) test
