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

| Platform | How-to |
| --- | --- |
| Linux | `curl https://linux.getgo.dev \| sh` |
| MacOS | `curl https://macos.getgo.dev \| sh` |
| Windows | Go to [the Releases page](./releases) |

All releases will also include binaries for all three supported packages with source code so you can build it yourself. See [the section on Compilation](#compilation) for details.



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
