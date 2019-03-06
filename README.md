# `GoDev`
GoDev is a live-reload development tool with first class support for Golang development.

- - -

## Getting Started
The following systems are somewhat supported:

- Linux
- MacOS
- Windows (*requires more testing, feel free to raise issues*)

You will also require **Go > 1.11.x** for GoDev to work out of the box.



### Installation
All releases will also include binaries for all three supported packages with source code so you can build it yourself. See [the section on Compilation](#compilation) for details.



#### Linux
Coming soon...



#### MacOS
Coming soon...



#### Windows
Coming soon...



### Usage: Develop with live-reload
Running `godev` without flags is the easiest way to get started with live-reload-enabled development:

```sh
godev
```

> GoDev runs `go mod vendor` to install dependencies, `go build -o bin/app` to build your application, and lastly it runs your app through `bin/app`. You might have to run `chmod +x bin/app` on the first build.



### Usage: Test with live-reload
To run the tests, simply specify the `--test` flag.

```sh
godev --test
```

To view verbose logs, append the `--vv` flag.

```sh
godev --test --vv
```



### Usage: Initialise a directory for Go development
To initialise a directory for development in the Golang language with Go Modules for package management, use the `--init` flag.

```sh
godev --init
```

> If you'd like to preview files before you install them, you can use the `--view` flag to check out the file first.



### Usage: Via Docker container
GoDev is also available as a Docker image at `zephinzer/godev:latest`. To get started using the Docker image, run:

```sh
# create the cache directory because volume definitions
# will cause a directory to be created as `root` which
# you will have trouble deleting on the host
mkdir -p ./.cache/pkg;

# instantiate the docker image into a container
docker run -it \
  -u $(id -u) \
  -v "$(pwd):/go/src/app" \
  -v "$(pwd)/.cache/pkg:/go/pkg" \
  zephinzer/godev:latest \
  godev
# ^ add any other flags you want after godev

```

- - -

## Advanced Usage
While GoDev was written focused on Golang development happiness, it can also be used for projects in other languages. Use the configuration flags to adjust it to your needs



### Flags

#### Run Modes
By default, GoDev will run for live-reload in development. This results in the default execution groups of:

1. `go mod vendor`
1. `go build -o ${BUILD_OUTPUT}` (*see `--output`*)
1. `${BUILD_OUTPUT}`

##### `--test`
Tells GoDev to run in test mode. This changes the default execution groups so that the following are run instead:

1. `go mod vendor`
1. `go build -o ${BUILD_OUTPUT}`  (*see `--output`*)
1. `go test ./... -coverprofile c.out`

##### `--init`
Specifying this flag triggers a directory initialisation flow which asks if you would like to initialise some files/directories if they are not found. These are:

1. Git repository (.git)
1. .gitignore
1. go.mod
1. main.go
1. Dockerfile
1. .dockerignore
1. Makefile 

##### `--view`
Specifying this flag with the name of a file prints the file to your terminal. For example, `godev --view main.go` will print the `main.go` file which `--init` will seed for you if you say yes.

#### Logs Verbosity

##### `--vv`
Defines verbose logs (debug level). Useful for debugging or if you'd like some insights into what triggered your job and to debug the pipeline for your specified execution groups.

##### `--vvv`
Defines super verbose logs (trace level). More useful if you're developing GoDev itself to trace the flow of events.

##### `--silent`
Tells GoDev to keep completely quiet. Only panic level logs are printed before GoDev exits with a non-zero status code.

#### Configuration

##### `--dir`
Specifies the directory for commands from GoDev to run from.

Default: Current working directory

##### `--watch`
Specifies the directory for GoDev to watch for changes recursively in.

Default: Current working directory

##### `--exec`
Specifies a single execution group. Commands specified in an execution group run in parallel.

Use multiple of these to define multiple execution groups. The execution groups run in sequence themselves.

Default (without `--test`):
1. `go mod vendor`
2. `go build -o ${BUILD_OUTPUT}` (*see `--output`*)
3. `${BUILD_OUTPUT}` (*see `--output`*)

Default (with `--test`):
1. `go mod vendor`
2. `go build -o ${BUILD_OUTPUT}`  (*see `--output`*)
3. `go test ./... -coverprofile c.out`

##### `--exec-delim`
Specifies the delimiter used in the `--exec` flag for separating commands. This flag finds its use if the command you wish to run contains a command as an argument.

Default: `,`

##### `--exts`
Defines a comma separated list of extensions (without the dot) to trigger a file system change event.

Default: `go,Makefile`

##### `--ignore`
Defines names of files/directories to ignore.

Default: `bin,vendor`

##### `--output`
Defines the path to the built output

Default: `bin/app`

#### `--rate`
Defines the rate at which file system change events are batched. Modifying this would be useful if you find that commands being run in your execution groups take longer than 2 seconds and modify files resulting in a never-ending file system change trigger loop.

Default: `2s`

##### `--version`
Prints the version of GoDev.

- - -

## Contributing

### Repository Setup
Run the following to clone this repository:

```sh
git clone git@github.com:zephinzer/godev.git;
```

Then copy the `sample.properties` into `Makefile.properties`:

```sh
cp sample.properties Makefile.properties;
```

### Dependency Installation
Dependencies are stored in the `./vendor` directory. Run the following to populate the dependencies:

```sh
make deps
```

### Static file generation
For static files that GoDev can initialise, a Go generator is used. The files can be found at `./data/generate` and the code to generate the file at `./data.go` can be found at `./data/generate.go`. The `./data.go` is generated with every build, but if you want to generate it manually, run:

```sh
make generate
```

### Development
Development is done in the `./dev` directory. Unfortunately, since this is a live-reload tool, there is no live-reload for the live-reload, so we have to re-run the application every time we make changes for them to be visible. The command to re-compile and re-run GoDev while working with `./dev` is:

```sh
make start
```

### Testing
To run the tests in watch mode:

```sh
make test
```

### Versioning
We try to follow [semver versioning]((https://semver.org/)) as far as possible. This means:

- Patch version bumps for bug fixes
- Minor version bumps for new flags/behaviours
- Major version bumps for deprecation of flags/behaviours

To get the current version:

```sh
make version.get
```

To bump the version:

```sh
# bump the patch version
make version.bump

# bump the minor version
make version.bump VERSION=minor

# bump the major version
make version.bump VERSION=major
```

### Compilation to Binary
To compile GoDev, run:

```sh
make compile
```

### Building the Docker Image
To build the GoDev image, run:

```sh
make docker
```

### Configuring the CI Pipeline
A Travis pipeline is included.

The following are variables that need to be defined for your pipeline:

| Environment Variable | Description |
| --- | --- |
| `DOCKER_IMAGE_REGISTRY` | Hostname of the Docker registry we are pushing to |
| `DOCKER_IMAGE_NAMESPACE` | Namespace of the Docker image (docker.io/THIS/image:tag) |
| `DOCKER_IMAGE_NAME` | Name of the Docker image (docker.io/namespace/THIS:tag |
| `DOCKER_REGISTRY_USERNAME` | Username for the Docker registry (when not specified, does not release to DockerHub) |
| `DOCKER_REGISTRY_PASSWORD` | Password for the Docker registry (when not specified, does not release to DockerHub) |
| `GITHUB_OAUTH_TOKEN` | GitHub personal access token for deploying binaries to the release page |
| `GITHUB_REPOSITORY_URL` | Clone URL of the GitHub repository (when not specified, does not release to GitHub) |
| `GITHUB_SSH_DEPLOY_KEY` | Base64 encoded private key that matches a public key listed in your Deploy Keys for the project. Run `make ssh.keys` to generate this. |

You will also need to go to your GitHub repository's **Settings > Deploy keys** and add the public key generated from `make ssh.keys` (the public key should be at `./bin/id_rsa.pub`, use the `./bin/id_rsa_b64` contents for the `GITHUB_SSH_DEPLOY_KEY` variable).

### TODOS
- Releasing to DockerHub
- Releasing to Brew
- Releasing to Chocolatey
- Releasing to GitHub
- Releasing to GitLab

- - -

## Architecture Notes

### Components
#### Watcher
- Watches the file system recursively at a directory level, watches new directories as they are created, sends notifications through a channel to the main process
- Batches file system changes and notifies the main process through a channel

#### Runner
- Handles the (re-)execution/termination of defined execution groups and commands
- Triggered through a function call that will terminate existing pipelines and restart them

#### Main Process
- Coordinates the batched file system changes from Watcher and triggers the Runner to start executing a pipeline

### Concepts
#### Pipeline
- Set of execution groups that run in sequence
- One pipeline per instantantiation of godev

#### Execution Groups
- Group of commands to run in parallel
- Execution groups run in sequence themselves

#### Command
- Atomic execution unit that runs a command using the user’s shell

- - -

## Support

### Work Hours
This is a side-project of mine meant to support my own development needs. I have a day job, so unless I have an urgent need while using this in my professional work, most of my code-level work on this repository will be done during weekends. Pull requests are however supported throughout the week!(:

Thanks for understanding!

### If You Really Like This
If you really like my work and would like to support me, you can find my Patreon at:



- - -

# Licensing
The binary and source code is licensed under the permissive MIT license. See [the LICENSE file](./LICENSE) for the full text.

- - -

# Cheers 😎
Leave me a 🌟 or watch this repository to indicate your interest in my sustained development on this. It'll help me decide whether or not I should deprecate this once my own use case for this is over.
