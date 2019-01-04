# Go Develop 🤓
A set of scripts/tools packaged in a Docker image to quickly get up-and-writing with Golang.

# Scope of Work
- ✅ Use Docker to provision an environment for a Golang 1.x application
- ✅ Live-reload of application using `go build`
- ✅ Live-updating of application's dependenices (`go mod` used here)
- ✅ Live-reload of tests using `go test`
- ✅ Statically link production binary using `go build`
- ✅ Docker image building using the `scratch` image
- ✅ Still be able to run things like `go (mod|build|run|test)` on the host system

# System Requirements
- [Make](https://www.gnu.org/software/make/)
- Docker ([RHEL](https://docs.docker.com/install/linux/docker-ce/centos/) [Fedora](https://docs.docker.com/install/linux/docker-ce/fedora/) [Ubuntu](https://docs.docker.com/install/linux/docker-ce/ubuntu/) [OS X](https://docs.docker.com/docker-for-mac/install/) [Windows](https://docs.docker.com/docker-for-windows/install/))
- [Docker Compose](https://docs.docker.com/compose/install/#install-compose)

> You don't need Go installed or a `GOPATH` for this 🥂

# Usage



## TL;DR
The quickest way to try this out is to copy the following into a Makefile in your target directory and run `make init`:

```makefile
GOLANG_DEV_VERSION=latest

# initialises this directory - use once only
init:
	@$(MAKE) _dev ARG="init"
# builds the application - outputs an `app` binary
build: 
	@$(MAKE) _dev ARG="build"
# runs tests in watch-mode
test: build
	@$(MAKE) _dev ARG="test"
# runs tests once
test.once: build
	@$(MAKE) _dev ARG="test -coverprofile c.out"
# runs the application on the host network
start: build
	@$(MAKE) _dev ARG="start"
# creates a shell in a fresh container generated from the image, usable for development on non-linux machines
shell:
	$(MAKE) dev ARG="shell"
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
    zephinzer/golang-dev:$(GOLANG_DEV_VERSION) ${ARG}
```

Check out [the documentation below on how to bundle your application into a `scratch` Docker image](#building-into-a-docker-image).



## Getting Started With Just The `golang-dev` Image
If you already have a Makefile/don't want a Makefile (but why?), running the following in any directory should get you up and running:

```sh
docker run -it \
  -u ${UID} \
  -v "$(pwd)/.cache/pkg:/go/pkg" \
  -v "$(pwd):/go/src/app" \
  zephinzer/golang-dev:latest init;
```

The script should ask a series of questions which if you agree to all, will create a Git repository, add an appropriate .gitignore file, initialise `go mod`, and provision a Dockerfile you can use to build your own image.

Notes:
- The `-u` parameter sets the guest user to use your host user ID so that new files aren't created as `root`, causing problems on your host machine when you attempt to delete them.
- The mapping of the `.cache` directory allows for faster subsequent dependency processing by `go mod`.



## Live-Reload, Live-Dependency-Update, Live-Testing Development
Assuming you've completed the above step, create a `main.go` if there isn't already one, and run the following to start development:

```sh
docker run -it \
  -u ${UID} \
  --network host \
  -v "$(pwd)/.cache/pkg:/go/pkg" \
  -v "$(pwd):/go/src/app" \
  zephinzer/golang-dev:latest start;
```

To start the tests, create a new terminal in the same directory and run:

```sh
docker run -it \
  -u ${UID} \
  --network host \
  -v "$(pwd)/.cache/pkg:/go/pkg" \
  -v "$(pwd):/go/src/app" \
  zephinzer/golang-dev:latest test;
```

**Notes**
- Making changes to `*.go` will result in a live-reload
- Add a new `import` will cause the script to start downloading the new dependency
- Removing that `import` will simply cause the dependency to be removed - simple stuff
- If you need a port open without using `--network host`, specify the flag `-p "HOSTPORT":"CONTAINERPORT"` in the `docker run` command
- If a new file is created, you'll need to send an interrupt to the process and start it again to begin watching it. I have no idea how to automatically watch new files but if you do, please submit a pull request. Thanks in advance!



## Building a Binary
Assuming you have a `main.go`, building a binary is as simple as running:

```sh
docker run -it \
  -u ${UID} \
  -v "$(pwd)/.cache/pkg:/go/pkg" \
  -v "$(pwd):/go/src/app" \
  zephinzer/golang-dev:latest build;
```

The binary will appear in your current directory and be named `app.$GOOS.$GOARCH`. It will be postfixd with a `.exe` if `$GOOS` specifies a Windows build.

**Notes**:
- The build is run with `CGO_ENABLED=0`, `GOOS=linux`, and `GOARCH=amd64`. To change this, specify your environment variables as part of the `docker` command. [Documentation can be found here under the `--env` or `--env-file` flag](https://docs.docker.com/engine/reference/commandline/run/#options).



## Building Into a Docker Image
If you denied the `init` script from creating a Dockerfile, you can run it again to get it. Otherwise, the Dockerfile script should be:

```Dockerfile
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

Copy the above and paste it in a file named Dockerfile in your root directory. After that, you can build your image by running:

```sh
docker build -t yourname/yourimage:yourtag .;
```



# Advanced Usage



## Usage in a `docker-compose.yml`
Add the following in your Docker Compose to use this as part of a Docker provisioned environment:

```yaml
version: "3.5" # or anything you're using
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



# Technical Footnotes



## Why `golang-dev` Came About
- **Golang requirements**. Go's infamous `$GOPATH` is a pain-in-the-filesystem. Developing via a Docker environment abstracts this away. Also, you no longer have to have Go installed on your system since managing different Go versions is another pain-in-the-filesystem.
- **Live-reload tooling**. [Realize](https://github.com/oxequa/realize) works for applications but does not run tests. Same for [Gin](https://github.com/codegangsta/gin) which seems . Only [GoConvey](https://github.com/smartystreets/goconvey) does it well, but it uses a browser interface instead of a CLI. Go Convey also only outputs a PASS or FAIL. I wanted to see error logs on FAIL which it didn't provide me.
- **Auto-download dependencies**. None of Realize, Gin and GoConvey automatically downloads dependencies either. While this is super useful, it also requires context of which package manager to use. This project enforces the use of `go mod` which accomodates previous package managers such as `glide` and `go dep`.



## Where `golang-dev` Fits In
I use this for development of information-systems type of software - basically, CRUD services. It might work for your own development, it might not. Drop me a pull request if there's something you can and are willing to fix/add, drop me an issue otherwise and I'll see if I can get it done!



- - -



# Development/Hacking
If you're interested in working on this, read on!

The contribution mechanism is roughly the same as other open-source projects:

1. Fork this repository
1. Make your changes on `master`
1. Make a pull request

## Code
The main logic of how this works is written as shell scripts [in the `/scripts` directory](./scripts).

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

This publishes two images - one with the version as recorded by the Git tags, another with the version of Golang, and a last one with the `latest` tag.



# Other Things

## Licensing
This project is licensed under the MIT license. See [the LICENSE file](./LICENSE) for the full text.

## All Relevant Links
- [GitHub Repository](https://github.com/zephinzer/golang-dev)
- [DockerHub Repository](https://hub.docker.com/r/zephinzer/golang-dev)

## Support/Work Hours
This is a side-project of mine meant to support my own development needs. I have a day job, so unless I have an urgent need while using this in my professional work, most of my code-level work on this repository will be done during weekends. Pull requests are supported throughout the week!(:

Thanks for understanding!

## Projects Using This
- [usvc/accounts](https://github.com/usvc/accounts)

## Cheers 😎
Leave me a 🌟 or watch this repository to indicate your interest in my sustained development on this. It'll help me decide whether or not I should deprecate this once my own use case for this is over.
