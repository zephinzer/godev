package main

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

// DefaultFileExtensions - default commma-separated list of file extensions to watch for
const DefaultFileExtensions = "go"

// DefaultIgnoredNames - default comma-separated list of file/dir names to ignore
const DefaultIgnoredNames = "bin,vendor"

// DefaultRefreshRate - default duration at which to handle file system events
const DefaultRefreshRate = 2 * time.Second

// Config configures the main application entrypoint
type Config struct {
	RunInit        bool
	RunTest        bool
	RunModWatch    bool
	FileExtensions commaSeparatedStringArray
	IgnoredNames   commaSeparatedStringArray
	Rate           time.Duration
	WatchDirectory string
}

// InitConfig creates a configuration from environment variables and flags
func InitConfig() *Config {
	currentWorkingDirectory := getCurrentWorkingDirectory()
	config := &Config{}
	flag.BoolVar(&config.RunInit, "init", false, "when this flag is specified, godev initiaises the current directory")
	flag.BoolVar(&config.RunModWatch, "watch", false, "when this flag is specified, godev runs the command in watch mode if applicable")
	flag.BoolVar(&config.RunTest, "test", false, "when this flag is specified, godev runs the tests with coverage")
	flag.Var(&config.FileExtensions, "exts", fmt.Sprintf("comma separated list of file extensions to watch (defaults to: %s)", DefaultFileExtensions))
	flag.Var(&config.IgnoredNames, "ignore", fmt.Sprintf("comma separated list of names to ignore (defaults to: %s)", DefaultIgnoredNames))
	flag.DurationVar(&config.Rate, "rate", DefaultRefreshRate, "specifies the refresh rate of the file system watch")
	flag.StringVar(&config.WatchDirectory, "dir", currentWorkingDirectory, "specifies the directory to watch")
	flag.Parse()
	if len(config.IgnoredNames) == 0 {
		config.IgnoredNames = strings.Split(DefaultIgnoredNames, ",")
	}
	if len(config.FileExtensions) == 0 {
		config.FileExtensions = strings.Split(DefaultFileExtensions, ",")
	}
	return config
}
