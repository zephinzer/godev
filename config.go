package main

import (
	"flag"
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
	ExecGroups        ConfigMultiflagString
	FileExtensions    ConfigCommaDelimitedString
	IgnoredNames      ConfigCommaDelimitedString
	LogLevel          LogLevel
	LogSilent         bool
	LogSuperVerbose   bool
	LogVerbose        bool
	Rate              time.Duration
	RunInit           bool
	RunTest           bool
	RunVersion        bool
	RunView           bool
	View              string
	WatchDirectory    string
	WorkDirectory     string
}

// InitConfig creates a configuration from environment variables and flags
func InitConfig() *Config {
	currentWorkingDirectory := getCurrentWorkingDirectory()
	config := &Config{}
	flag.StringVar(&config.BuildOutput, "output", DefaultBuildOutput,
		"specifies the path to the built binary relative to the watch directory (applicable only when --exec is not specified)")
	flag.StringVar(&config.CommandsDelimiter, "exec-delim", DefaultCommandsDelimiter,
		"delimiter character to use to split commands within an execution group (useful if your commands themselves contain commas)")
	flag.Var(&config.ExecGroups, "exec",
		"list of comma-separated commands to run (specify multiple --execs to define more execution groups)")
	flag.Var(&config.FileExtensions, "exts",
		fmt.Sprintf("comma separated list of file extensions to watch (defaults to: %s)", DefaultFileExtensions))
	flag.Var(&config.IgnoredNames, "ignore",
		fmt.Sprintf("comma separated list of names to ignore (defaults to: %s)", DefaultIgnoredNames))
	flag.BoolVar(&config.LogSilent, "silent", false,
		"it is better to remain silent at the risk of being thought a fool, than to talk and remove all doubt of it")
	flag.BoolVar(&config.LogSuperVerbose, "vvv", false,
		"show me everything that has been logged")
	flag.BoolVar(&config.LogVerbose, "vv", false,
		"show sensibly verbose logs")
	flag.DurationVar(&config.Rate, "rate", DefaultRefreshRate,
		"specifies the duration between two bathced file system event triggers (increase this to the duration which commands that modify files take to complete)")
	flag.BoolVar(&config.RunInit, "init", false,
		"when this flag is specified, godev runs an initialisation procedure in the current directory or the directory specified by --dir")
	flag.BoolVar(&config.RunTest, "test", false,
		"when this flag is specified, godev runs the tests with coverage")
	flag.BoolVar(&config.RunVersion, "version", false,
		"display the version number")
	flag.StringVar(&config.View, "view", "",
		"check out the original content of a file that godev provisions when --init is specified")
	flag.StringVar(&config.WatchDirectory, "watch", currentWorkingDirectory,
		"specifies the directory to watch")
	flag.StringVar(&config.WorkDirectory, "dir", currentWorkingDirectory,
		"specifies the directory to operate in")
	flag.Parse()
	config.assignDefaults()
	config.interpretLogLevel()
	return config
}

func (config *Config) interpretLogLevel() {
	if config.LogVerbose {
		config.LogLevel = "debug"
	}
	if config.LogSuperVerbose {
		config.LogLevel = "trace"
	}
	if config.LogSilent || config.RunVersion {
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
