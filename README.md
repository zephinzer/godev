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
Specifies the directory to run `godev` from.

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


#### Meta-data

##### `--version`



### JavaScript Project

- - -

## Contributing
- Cloning
- Dependency Installation
- Static file generation
- Development
- Testing
- Versioning
- Compilation
- Releasing to DockerHub
- Releasing to Brew
- Releasing to Chocolatey
- Releasing to GitHub
- Releasing to GitLab

- - -

## Architecture Notes

### Watcher
- Watches the file system recursively at a directory level, watches new directories as they are created, sends notifications through a channel to the main process
- Batches file system changes and notifies the main process through a channel
### Runner
- Handles the (re-)execution/termination of defined execution groups and commands
- Triggered through a function call that will terminate existing pipelines and restart them
### Main Process
- Coordinates the batched file system changes from Watcher and triggers the Runner to start executing a pipeline
### Concept: Pipeline
- Set of execution groups that run in sequence
- One pipeline per instantantiation of godev
### Concept: Execution Groups
- Group of commands to run in parallel
- Execution groups run in sequence themselves
### Concept: Command
- Atomic execution unit that runs a command using the user’s shell

- - -

## Support
- Weekends only - feel free to raise a PR
- Patreon Link

- - -

## Licensing
- MIT

# Cheers
