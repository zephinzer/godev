package main

import (
	"os"

	"github.com/urfave/cli"
)

func initCLI() *CLI {
	app := &CLI{}
	app.config = &Config{}
	app.logger = InitLogger(&LoggerConfig{
		Name:   "cli",
		Format: "production",
		Level:  "trace",
	})
	app.rawLogger = InitLogger(&LoggerConfig{
		Name:   "cli",
		Format: "raw",
		Level:  "trace",
	})
	instance := cli.NewApp()
	instance.Name = "godev"
	instance.Usage = "a development tool for golang"
	instance.Description = "golang development tool with project bootstrap, live-reload, and auto-dependency retrieval powers"
	instance.Version = Version
	instance.Action = getDefaultAction(app.config)
	instance.Commands = []cli.Command{
		getInitCommand(app.config),
		getTestCommand(app.config),
		getVersionCommand(app.config, app.rawLogger),
		getViewCommand(app.config, app.rawLogger),
	}
	instance.Flags = getDefaultFlags()
	app.instance = instance
	return app
}

type CLI struct {
	config    *Config
	instance  *cli.App
	logger    *Logger
	rawLogger *Logger
}

func (app *CLI) Start(args []string, after func(*Config)) {
	app.instance.After = func(c *cli.Context) error {
		after(app.config)
		return nil
	}
	if err := app.instance.Run(args); err != nil {
		app.logger.Error(err)
		app.logger.Warn("exiting with status code 1")
		os.Exit(1)
	}
}
