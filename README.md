# Golang Development Image
This repository contains a Docker image used for an easier Golang 1.x development experience

# Requirements
You'll need:

1. Docker
2. A sense of adventure

# Opinions forced on you

1. You should be using Git
1. Usage of `go mod` is the only way to manage dependencies
1. Having a vendor library is good for your health
1. Live-reload is happiness on a CLI
1. Live-reload should also live-install dependencies added to your `*.go`

# Intentions

1. Meant for information systems development
1. Host system independent methodology
1. Destined for containerisation

# Usage

## Directory Initialisation
Run the following in your directory of choice:

```sh
docker run -it -u ${UID} -v "$(pwd):/go/src/app" zephinzer/golang-dev:latest init;
```

## Project Setup
Setup a `Dockerfile` in the directory of your choice with:

```
FROM zephinzer/golang-dev:latest as development
COPY . /go/src/app
ENTRYPOINT [ "start" ]

FROM development as build
RUN build
ENTRYPOINT [ "/go/src/app/app" ]

FROM scratch as production
COPY --from=build /go/src/app/app /
ENTRYPOINT [ "/app" ]
```

Add a `main.go` to the same directory:

```go
package main

import (
  "fmt"
)

func main() {
  fmt.Println("hello world")
}
```

## Developing with Live-Reload
Run the following to start your application in live-reload:

```sh
docker run -it \
  -u ${UID} \
  -v "$(pwd):/go/src/app" \
  zephinzer/golang-dev:latest start;
```

If you need to expose ports (3000 in the example):

```sh
docker run -it \
  -u ${UID} \
  -v "$(pwd):/go/src/app" \
  -p "3000:3000" \
  zephinzer/golang-dev:latest start;
```

## Testing with Live-Reload
Run the following to test your application in watch mode:

```sh
docker run -it \
  -u ${UID} \
  -v "$(pwd):/go/src/app" \
  zephinzer/golang-dev:latest \
  test;
```

## Compilation to Binary
Run the following to compile your application:

```sh
docker run \
  -u ${UID} \
  -v "$(pwd):/go/src/app" \
  zephinzer/golang-dev:latest \
  build;
```

You should now see a binary named `app` in that directory which you can redistribute.

## Docker Image Packaging
Alternatively, to package it into a nice Docker image:

```sh
docker build \
  -t yourusername/imagename:latest \
  .
```

Run `docker images` and you should see your image. This image is built from `scratch`, which is super minimal.

Try to run it:

```sh
docker run yourusername/imagename:latest;
```

## TL;DR (*Gimme a Makefile*)

Pretty self-explanatory, copy this into a Makefile in an empty directory and run `make init` to get started!

```makefile
GOLANG_DEV_VERSION=latest

init: # initialises this directory - use once only
	@$(MAKE) _dev ARG="init"
build: # builds the application - outputs an `app` binary
	@$(MAKE) _dev ARG="build"
test: build # runs tests in watch-mode
	@$(MAKE) _dev ARG="test"
test.once: build # runs tests once
	@$(MAKE) _dev ARG="test -coverprofile c.out"
start: build # runs the application on the host network
	@$(MAKE) _dev ARG="start"
shell: # creates a shell in a fresh container generated from the image, usable for development on non-linux machines
	$(MAKE) dev ARG="shell"
version.get: # retrieves the latest version we are at
	@docker run -v "$(CURDIR):/app" zephinzer/vtscripts:latest get-latest -q
version.bump: # bumps the version by 1: specify VERSION as "patch", "minor", or "major", to be specific about things
	@docker run -v "$(CURDIR):/app" zephinzer/vtscripts:latest iterate ${VERSION} -i
_dev: # base command to run (do not use)
	@docker run \
    -it \
    --network host \
    -u $$(id -u) \
    -v "$(CURDIR)/.cache/pkg:/go/pkg" \
    -v "$(CURDIR):/go/src/app" \
    zephinzer/golang-dev:$(GOLANG_DEV_VERSION) ${ARG}
```

# Advanced Usage

## Running on Host Network

```sh
docker run \
  -u ${UID} \
  --network host \
  -v "$(pwd):/go/src/app" \
  zephinzer/golang-dev:latest \
  start;
```

## Running within Docker Compose

In Docker Compose, add the following:

```yaml
version: "3.5"
services:
  # ...
  app:
    image: zephinzer/golang-dev:latest
    ports: # if needed
    - "3000:3000"
    user: "${UID}"
    entrypoint: ["start"]
    volumes:
    - "./.cache/pkg:/go/pkg"
    - "./:/go/src/app"
    # ...
  # ...
```

In Makefile, add a `start.once` recipe (if needed), and change `start` recipe to:
```makefile
# ...
start: # starts the development environment
	@UID=$$(id -u) docker-compose up -d ${ARGS} app
start.once: build # runs the application on the host network
	@$(MAKE) _dev ARG="start"
# optional: if you want an easy way to get the logs:
logs: # displays the application logs
	@docker logs -f $$(docker ps | grep $$(basename $$(pwd)) | grep app | cut -f 1 -d ' ')
# ...
```

## Building Binaries for Other Architectures

An example for Windows follows:

```sh
docker run \
  -u ${UID} \
  -v "$(pwd):/go/src/app" \
  --env GOARCH=amd64 \
  --env GOOS=windows \
  zephinzer/golang-dev:latest \
  build;
```

[Check out this page](https://golang.org/doc/install/source#environment) for all possible `GOARCH` and `GOOS`es.

# Development/Hacking

## Code
The main logic of how this works is written in `bash` [in the `/scripts` directory](./scripts).

The `Dockerfile` simply copies [the `/scripts`](./scripts) in and adds it to the `$PATH`.

## Testing
Tests are contained [in the `./test` directory](./test) but you can run it from the root using `make test`.

## Building
To build the Docker image, run `make build`.

## Versioning
To bump a patch version, run `make version.bump`.

To bump a minor versoin, run `make version.bump VERSION=minor`.

To bump a major versoin, run `make version.bump VERSION=major`.

## Publishing
To publish the Docker image, run `make publish`.

This publishes two images - one with the version as recorded by the Git tags, another with the version of Golang.

# License

This project is licensed under the MIT license. See [the LICENSE file](./LICENSE) for the full text.
