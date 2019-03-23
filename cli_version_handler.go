package main

import (
	"github.com/urfave/cli"
)

func getVersionCommand(config *Config, logger *Logger) cli.Command {
	return cli.Command{
		Action:      getVersionAction(config, logger),
		Aliases:     []string{"v"},
		Description: "print the version",
		Flags:       getVersionFlags(),
		Name:        "version",
		Usage:       "print the version",
	}
}

func getVersionFlags() []cli.Flag {
	return []cli.Flag{
		getFlagSemver(),
		getFlagCommit(),
	}
}

func getVersionAction(config *Config, logger *Logger) cli.ActionFunc {
	return func(c *cli.Context) error {
		config.RunVersion = true
		config.interpretLogLevel()
		if c.Bool("semver") {
			logger.Info(Version)
		} else if c.Bool("commit") {
			logger.Info(Commit)
		} else {
			logger.Infof("godev %s-%s", Version, Commit)
		}
		return nil
	}
}
