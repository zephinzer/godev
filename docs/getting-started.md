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

For further configuration, see the [Configuration Flags](./configuration-flags.md) page.

- - -

[Back](../)
