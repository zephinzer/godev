# Getting Started with `godev`

## Installation

### Go Get

Use `go get` to install this package locally:

```sh
go get github.com/zephinzer/godev
```

### Linux

Coming soon.

### MacOS

Coming soon.

### Windows

Coming soon.

## Usage

### Basic/Quick Start

To run your Go application with live-reload:

```sh
godev
```

> This will make `godev` listen for file changes and on every file change cycle, perform: 1) dependency installation via `go mod vendor` 2) a `go build` 3) a run of your built application

To run your Go tests with live-reload:

```sh
godev --test
```

> This will make `godev` listen for file changes and on every file change cycle, perform: 1) dependency installation via `go mod vendor` 2) a `go build` 3) a run of your tests using `go test ./...`

To define a different directory to watch:

```sh
godev --watch /path/to/else/where
```

> `godev` will watch a different directory from the current working directory

To define extensions to watch:

```sh
godev --exts go,Makefile,Dockerfile
```

> Will watch for changes in \*.go, (\*.)Makefile, or (*.)Dockerfile files.

### Advanced

#### Custom Commands
To run a set of custom commands, use the `--exec` flag. Each `--exec` defines an execution group where all commands are run in parallel. Multiple `--exec`s form the pipeline where each execution group will run in sequence.

For example, given:

```sh
godev --exec 'echo 1,echo 2' --exec 'echo 3'
```

The `echo 1` and `echo 2` will be run in parallel first, after which, `echo 3` will be executed.

If your command uses a comma, specify the `--exec-delim` flag to change the delimiter between commands:

```sh
godev --exec 'echo 1|echo 2' --exec 'echo 3' --exec-delim '|'
```
