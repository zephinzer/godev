package main

import (
	"fmt"

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
		Description: "bootstrap the current working directory for Golang development",
		Name:        "init",
		Usage:       "bootstrap the current working directory for Golang development",
	}
}

func getCommandTest(config *Config) cli.Command {
	return cli.Command{
		Action:      getActionTest(config),
		Aliases:     []string{"t"},
		Description: "run tests in live-reload mode",
		Flags:       getCommandTestFlags(),
		Name:        "test",
		Usage:       "run tests in live-reload mode",
	}
}

func getCommandVersion(config *Config) cli.Command {
	return cli.Command{
		Action:      getActionVersion(config),
		Aliases:     []string{"v", "ver"},
		Description: "print the version",
		Flags:       getCommandVersionFlags(),
		Name:        "version",
		Usage:       "print the version",
	}
}

func getCommandView(config *Config) cli.Command {
	description := "checkout files seeded with the init sub-command where the [filename] argument is one of:"
	for filename := range InitFileMap {
		description = fmt.Sprintf("%s\n  - %s", description, filename)
	}

	return cli.Command{
		Action:      getActionView(config),
		Aliases:     []string{"V"},
		ArgsUsage:   "[filename]",
		Description: description,
		Name:        "view",
		Usage:       "checkout files seeded with the init sub-command",
	}
}
