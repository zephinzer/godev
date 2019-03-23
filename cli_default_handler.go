package main

import (
	"strings"

	"github.com/urfave/cli"
)

func getDefaultFlags() []cli.Flag {
	return []cli.Flag{
		getFlagBuildOutput(),
		getFlagCommandsDelimiter(),
		getFlagEnvVars(),
		getFlagExecGroups(),
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

func getDefaultAction(config *Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		config.RunDefault = true
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
