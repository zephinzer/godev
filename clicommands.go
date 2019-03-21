package main

import (
	"github.com/urfave/cli"
)

func getCommands(config *Config) []cli.Command {
	return []cli.Command{
		getCommandInit(config),
		getCommandTest(config),
		getCommandVersion(config),
		getCommandView(config),
	}
}

func getCommandInit(config *Config) cli.Command {
	return cli.Command{
		Action:      getActionInit(config),
		Aliases:     []string{"i"},
		Name:        "init",
		Usage:       "bootstrap the current working directory for Golang development",
		Description: "bootstrap the current working directory for Golang development",
	}
}

func getCommandTest(config *Config) cli.Command {
	return cli.Command{
		Action:      getActionTest(config),
		Aliases:     []string{"t"},
		Flags:       getCommandTestFlags(),
		Name:        "test",
		Usage:       "run tests in live-reload mode",
		Description: "run tests in live-reload mode",
	}
}

func getCommandVersion(config *Config) cli.Command {
	return cli.Command{
		Action:      getActionVersion(config),
		Aliases:     []string{"v", "ver"},
		Flags:       getCommandVersionFlags(),
		Name:        "version",
		Usage:       "print the version",
		Description: "print the version",
	}
}

func getCommandView(config *Config) cli.Command {
	return cli.Command{
		Action:      getActionView(config),
		Aliases:     []string{"V"},
		Name:        "view",
		Usage:       "checkout files seeded with the init sub-command",
		Description: "checkout files seeded with the init sub-command",
	}
}
