package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func getInitCommand(config *Config) cli.Command {
	return cli.Command{
		Action:      getInitAction(config),
		Aliases:     []string{"i"},
		Description: "bootstrap the current working directory for Golang development",
		Flags:       getInitFlags(),
		Name:        "init",
		Usage:       "bootstrap the current working directory for Golang development",
	}
}

func getInitFlags() []cli.Flag {
	return []cli.Flag{
		getFlagWorkDirectory(),
	}
}

func getInitAction(config *Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		config.RunInit = true
		config.WorkDirectory = c.String("dir")
		fmt.Println(config.WorkDirectory)
		config.assignDefaults()
		config.interpretLogLevel()
		return nil
	}
}
