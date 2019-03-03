# `godev` Configuration Flags

Following are flags you can use with `godev` sorted in alphabetical order:

## `--exec`
Specifies an execution group. Use this flag multiple times to specify multiple execution groups. Each execution group runs in parallel. All execution groups run in sequence.

> Defaults to running `go mod vendor`, `go build ...`, and `/bin/app`. If `--test` is specified, defaults to running `go mod vendor`, `go build ...`, and `go test ./...`.

## `--exts`
Sets the watched extensions.

Usage #1: `godev --exts go,Makefile,Dockerfile`

Usage #2: `godev --exts js,ts,json`

> Defaults to \*.go and  (\*.)Makefile

## `--ignore`
Specifies names of files/directories to ignore when watching the current directory.

Usage: `--ignore bin,vendor`

> Defaults to ignoring `bin` and `vendor`.

## `--output`
Defines the relative path to the output binary. If you've specified a custom set of `--exec`s, this flag has no effect.

Usage: `--output bin/myapp`

> Defaults to `./bin/app` relative to your watched directory.

## `--rate`
Defines the duration of file system events batching.

Usage: `--rate 2s`

> Defaults to 2 seconds

## `--silent`
Disables all logs (panic level).

## `--test`
Indicates to run in test mode.

Note: If you've specified a custom set of `--exec`s, this flag has no effect.

Usage `--test`

## `--verison`
Outputs the version.

## `--view`
Previews a file that `godev` can initialise for you. Available options are: `dockerfile`, `makefile`, `go.mod`, `main.go`, `.gitignore`, `.dockerignore`.

Usage: `--view dockerfile`

## `--vv`
Turns on verbose logs (debug level).

Note: When used with `--test`, results in verbose test results.

## `--vvv`
Turns on super verbose logs (trace level).

Note: When used with `--test`, results in verbose test results.

## `--watch`
Sets the absolute path to be watched.

Usage: `godev --watch $(pwd)/bin`

> When this flag is not set, the watched directory defaults to the current working directory you ran `godev` from.
