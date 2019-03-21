package main

import (
	"os"

	"github.com/urfave/cli"
)

func initCLI() *cliApp {
	app := &cliApp{}
	app.config = &Config{}
	instance := cli.NewApp()
	instance.Name = "godev"
	instance.Usage = "a development tool for golang"
	instance.Description = "golang development tool with project bootstrap, live-reload, and auto-dependency retrieval powers"
	instance.Version = Version
	instance.Action = getActionDefault(app.config)
	instance.Commands = getCommands(app.config)
	instance.Flags = getGlobalFlags()
	app.instance = instance
	app.logger = InitLogger(&LoggerConfig{
		Name:   "cli",
		Format: "production",
		Level:  "trace",
	})
	return app
}

type cliApp struct {
	config   *Config
	instance *cli.App
	logger   *Logger
}

func (app *cliApp) Start(after func(*Config)) {
	app.instance.After = func(c *cli.Context) error {
		after(app.config)
		return nil
	}
	if err := app.instance.Run(os.Args); err != nil {
		app.logger.Panic(err)
	}
}
