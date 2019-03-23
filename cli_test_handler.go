package main

import (
	"strings"

	"github.com/urfave/cli"
)

func getTestCommand(config *Config) cli.Command {
	return cli.Command{
		Action:      getTestAction(config),
		Aliases:     []string{"t"},
		Description: "run tests in live-reload mode",
		Flags:       getTestFlags(),
		Name:        "test",
		Usage:       "run tests in live-reload mode",
	}
}

func getTestFlags() []cli.Flag {
	return []cli.Flag{
		getFlagBuildOutput(),
		getFlagCommandsDelimiter(),
		getFlagEnvVars(),
		getFlagFileExtensions(),
		getFlagIgnoredNames(),
		getFlagRate(),
		getFlagSilent(),
		getFlagSuperVerboseLogs(),
		getFlagVerboseLogs(),
		getFlagWatchDirectory(),
		getFlagWorkDirectory(),
	}
}

func getTestAction(config *Config) cli.ActionFunc {
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
