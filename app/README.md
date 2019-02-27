# Go Develop 🤓💾
A binary inspired by [the set of Makefile scripts](../) to quickly get up-and-writing with Golang.

# System Support
- Go 1.11.x
- Linux
- MacOS
- Windows (*not officially supported*)

# Usage

## Quick Start

### Installation

`TODO...`

- [ ] Linux: via binary download
- [ ] Linux: via cURL installation
- [ ] Linux: via `apt`
- [ ] MacOS: via Brew
- [ ] Windows: via Chocolatey

### Initialisation
Initialise a directory for Golang development using:

```sh
godev --init
```

### Development with Live-Re-Load/Build/Dependency download
To run the application in live-reload mode:

```sh
godev
```

The tool should run the following for you on every file system event related to `.go` files:

1. `go mod vendor` (retrieve dependencies)
1. `go build` (builds the binary)
1. `/path/to/binary` (runs the binary)

On a file system event, the tool should send a `SIGINT` first, wait for 5 seconds, send a `SIGTERM`, wait for another 5 seconds, and finally a `SIGKILL`.

### Testing with Live-Reload
To run the tests in live-reload mode, use:

```sh
godev --test
```

The tool should run the following for you on every file system event related to `.go` files:

1. `go mod vendor` (see above)
1. `go test ./...` (runs all files with file names ending in `_test.go` recursively)

## Execution Flags
Prepend a `-` before the flag when calling the `godev` command.

| Flag | Parameters | Description | Example values |
| --- | --- | --- |
| `exec` | `string...` | Comma separated list of commands (with arguments) to run in parallel. Each of this flag defines a set of commands to run in parallel (an execution group). When this is specified, `godev` stops being a Golang development tool and becomes a generic development tool for watching for file changes and executing shell scripts on changes | `go build,golint` |
| `exec-delim` | `string` | Delimiter for the commands specified in `--exec`. Defaults to a comma. | `/path/to/your/project` |
| `exts` | `string` | Comma separated strings denoting what extensions to look out for. Extensions should not contain the initial dot | `go,Makefile` |
| `help` | - | Displays the help message | - |
| `ignore` | `string` | Comma separated strings indicating what file or directory names to ignore | `vendor,cache` |
| `init` | - | Initialises the current directory for working with Golang | - |
| `output` | `string` | Relative path from the watch directory (`--dir`) to the binary **(only applicable when `--exec` is not defined)**  | `bin/app` |
| `rate` | `duration` | Time interval between when events are deduped and reported to the main handler. Defaults to 2 seconds. | `2s` |
| `test` | - | Runs the tests instead of the app | - |
| `view` | `string` | Checks out the file contents of the specified file. Can be any of `{dockerfile, makefile, .gitignore, .dockerignore, go.mod, main.go}` | `dockerfile` |
| `vv` | - | Logs in verbose mode (debug level) | - |
| `vvv` | - | Logs in super verbose mode (trace level) | - |
| `watch` | - | Absolute path to directory to consider the working directory. Defaults to the directory that `godev` is called from | `/path/to/your/project` |

# Development

## Architecture Overview
The application consists of two major components, the Watcher and the Runner. The Watcher watches for file system changes and notifies the main process. The Runner manages the running processes and is operated by the main process thread which restarts process pipelines when triggered by the Watcher's file system notifications.

## Getting Started

### Dependency Installation
Go Modules are used. To start development, install the dependencies with:

```sh
make install.deps
```

### Writing Code
The live-reloading from the Makefile version of this tool is used.

For convenience, flags have been set which make development easier:

```sh
make start
```

To run custom flags, use:

```sh
make start.prd ARGS="<__YOUR_FLAGS__>"
```

To run it in test mode, use:

```sh
make run ARGS="--test --ignore .cache,vendor,bin,data"
```

### Running Tests
We use standard Go tooling for this. To run the tests in live-reload mode, use:

```sh
make test
```

To run it once to generate coverage information:

```sh
make test.once
```

### Modifying the `--init` files
The seed files used during `--init` can be found at [data/generate](data/generate). After modification, use `go generate` to re-generate the `data.go` in this directory. You can use the convenience recipe:

```sh
make generate
```

### Building the Binary
The binary build is done with static linking. To build the binary and output it to `bin/godev`, use:

```sh
make compile
```

### Building the Docker Image

`TODO`

### Versioning Matters (For collaborators)
We try to stick to [Semver rules](https://semver.org/). Patch versions are for hot/quick-fixes that should never impact existing users via automated updates. Minor versions are for new functionalities. Major versions are for non-backward compatible functionalities that may break existing scripts.

Before releasing, use the following to bump the patch version:

```sh
make version.bump
```

And for major releases:

```sh
make version.bump VERSION=major
```

For minor, substitute the `"major"` with `"minor"`.



# Other Things

## Licensing
The binary, like [the Makefile](../), is licensed under the permissive MIT license. See [the LICENSE file](./LICENSE) for the full text.

## All Relevant Links

`TODO`

## Support/Work Hours
This is a side-project of mine meant to support my own development needs. I have a day job, so unless I have an urgent need while using this in my professional work, most of my code-level work on this repository will be done during weekends. Pull requests are supported throughout the week!(:

Thanks for understanding!

## Cheers 😎
Leave me a 🌟 or watch this repository to indicate your interest in my sustained development on this. It'll help me decide whether or not I should deprecate this once my own use case for this is over.
