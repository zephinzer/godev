package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

func getActionDefault(config *Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		config.BuildOutput = c.String("output")
		config.CommandsDelimiter = c.String("exec-delim")
		config.EnvVars = c.StringSlice("env")
		config.ExecGroups = c.StringSlice("exec")
		config.FileExtensions = strings.Split(c.String("exts"), ",")
		config.IgnoredNames = strings.Split(c.String("ignore"), ",")
		config.Rate = c.Duration("rate")
		config.WatchDirectory = c.String("watch")
		config.WorkDirectory = c.String("dir")
		config.assignDefaults()
		config.LogSilent = c.Bool("silent")
		config.LogVerbose = c.Bool("verbose")
		config.LogSuperVerbose = c.Bool("vverbose")
		config.interpretLogLevel()
		return nil
	}
}

func getActionInit(config *Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		config.RunInit = true
		return nil
	}
}

func getActionTest(config *Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		config.RunTest = true
		config.BuildOutput = c.String("output")
		config.CommandsDelimiter = c.String("exec-delim")
		config.EnvVars = c.StringSlice("env")
		config.FileExtensions = strings.Split(c.String("exts"), ",")
		config.IgnoredNames = strings.Split(c.String("ignore"), ",")
		config.Rate = c.Duration("rate")
		config.WatchDirectory = c.String("watch")
		config.WorkDirectory = c.String("dir")
		config.assignDefaults()
		config.LogSilent = c.Bool("silent")
		config.LogVerbose = c.Bool("verbose")
		config.LogSuperVerbose = c.Bool("vverbose")
		config.interpretLogLevel()
		return nil
	}
}

func getActionVersion(config *Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		config.RunVersion = true
		if c.Bool("semver") {
			fmt.Println(Version)
		} else if c.Bool("commit") {
			fmt.Println(Commit)
		} else {
			fmt.Printf("godev %s-%s", Version, Commit)
		}
		return nil
	}
}

func getActionView(config *Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		config.RunView = true
		return nil
	}
}
