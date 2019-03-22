package main

import (
	"fmt"
	"path"
	"strings"
	"time"
)

// DefaultBuildOutput - default relative path to watch directory to place built binaries in
const DefaultBuildOutput = "bin/app"

// DefaultCommandsDelimiter - default string to split --execs into commands with
const DefaultCommandsDelimiter = ","

// DefaultExecutionGroupsBase - default commands to run when no --execs are specified
var DefaultExecutionGroupsBase = []string{"go mod vendor"}

// DefaultFileExtensions - default commma-separated list of file extensions to watch for
const DefaultFileExtensions = "go,Makefile"

// DefaultIgnoredNames - default comma-separated list of file/dir names to ignore
const DefaultIgnoredNames = "bin,vendor"

// DefaultLogLevel - default log level from 'trace', 'debug', 'info', 'warn', 'error', 'panic'
const DefaultLogLevel = "info"

// DefaultRefreshRate - default duration at which to handle file system events
const DefaultRefreshRate = 2 * time.Second

// Config configures the main application entrypoint
type Config struct {
	BuildOutput       string
	CommandsDelimiter string
	EnvVars           ConfigMultiflagString
	ExecGroups        ConfigMultiflagString
	FileExtensions    ConfigCommaDelimitedString
	IgnoredNames      ConfigCommaDelimitedString
	LogLevel          LogLevel
	LogSilent         bool
	LogSuperVerbose   bool
	LogVerbose        bool
	Rate              time.Duration
	RunDefault        bool
	RunInit           bool
	RunTest           bool
	RunVersion        bool
	RunView           bool
	View              string
	WatchDirectory    string
	WorkDirectory     string
}

func (config *Config) interpretLogLevel() {
	if config.LogVerbose {
		config.LogLevel = "debug"
	}
	if config.LogSuperVerbose {
		config.LogLevel = "trace"
	}
	if config.LogSilent || config.RunVersion || config.RunView {
		config.LogLevel = "panic"
	}
}

func (config *Config) assignDefaults() {
	config.LogLevel = DefaultLogLevel
	config.BuildOutput = path.Join(config.WorkDirectory, "/"+config.BuildOutput)
	config.RunView = len(config.View) > 0
	if len(config.IgnoredNames) == 0 {
		config.IgnoredNames = strings.Split(DefaultIgnoredNames, ",")
	}
	if len(config.FileExtensions) == 0 {
		config.FileExtensions = strings.Split(DefaultFileExtensions, ",")
	}
	if len(config.ExecGroups) == 0 {
		if config.RunTest {
			testFlags := "-coverprofile c.out"
			if config.LogVerbose || config.LogSuperVerbose {
				testFlags = fmt.Sprintf("-v %s", testFlags)
			}
			config.ExecGroups = append(
				DefaultExecutionGroupsBase,
				fmt.Sprintf("go build -o %s", config.BuildOutput),
				fmt.Sprintf("go test ./... %s", testFlags),
			)
		} else {
			config.ExecGroups = append(
				DefaultExecutionGroupsBase,
				fmt.Sprintf("go build -o %s", config.BuildOutput),
				config.BuildOutput,
			)
		}
	}
}
