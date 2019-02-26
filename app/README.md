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

`TODO`

### Initialisation

```sh
godev --init
```

### Development with Live-Re-Load/Build/Dependency download

```sh
godev --watch
```

### Testing with Live-Reload

```sh
godev --test --watch
```

## Execution Flags

| Flag | Description | Example values |
| --- | --- | --- |
| `--help` | Displays the help message | - |
| `--init` | Initialises the current directory for working with Golang | - |
| `--test` | Runs the tests instead of the app | - |
| `--watch` | Runs the application/tests in watch mode | - |
| `--exts` | Comma separated strings denoting what extensions to look out for. Extensions should not contain the initial dot | `go,Makefile` |
| `--ignore` | Comma separated strings indicating what file or directory names to ignore | `vendor,cache` |
| `--rate` | Time interval between when events are deduped and reported to the main handler. Defaults to 2 seconds. | `2s` |
| `--dir` | Absolute path to directory to consider the working directory. Defaults to the directory that `godev` is called from | `/path/to/your/project` |
| `--output` | Relative path from the watch directory (`--dir`) to the binary **(only applicable when `--exec` is not defined)**  | `bin/app` |
| `--exec` | Comma separated list of commands (with arguments) to run in parallel. Each of this flag defines a set of commands to run in parallel (an execution group). When this is specified, `godev` stops being a Golang development tool and becomes a generic development tool for watching for file changes and executing shell scripts on changes | `go build,golint` |
| `--exec-delim` | Delimiter for the commands specified in `--exec`. Defaults to a comma. | `/path/to/your/project` |

# Development

## Architecture
The application consists of two major components, the Watcher and the Runner. The Watcher watches for file system changes and notifies the main process. The Runner manages the running processes and is operated by the main process thread which restarts process pipelines when triggered by the Watcher's file system notifications.
