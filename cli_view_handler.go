package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

type View struct {
	Filename string
	Data     string
}

var ViewMap = map[string]*View{
	"dockerfile":    &View{"Dockerfile", DataDockerfile},
	"makefile":      &View{"Makefile", DataMakefile},
	".dockerignore": &View{".dockerignore", DataDotDockerignore},
	".gitignore":    &View{".gitignore", DataDotGitignore},
	"main.go":       &View{"main.go", DataMainDotgo},
	"go.mod":        &View{"go.mod", DataGoDotMod},
}

func getViewCommand(config *Config, logger *Logger) cli.Command {
	description := "checkout files seeded with the init sub-command where the [filename] argument is one of:"
	for filename := range ViewMap {
		description = fmt.Sprintf("%s\n  - %s", description, filename)
	}

	return cli.Command{
		Action:      getViewAction(config, logger),
		Aliases:     []string{"V"},
		ArgsUsage:   "[filename]",
		Description: description,
		Name:        "view",
		Usage:       "checkout files seeded with the init sub-command",
	}
}

func getViewAction(config *Config, logger *Logger) cli.ActionFunc {
	return func(c *cli.Context) error {
		config.RunView = true
		config.View = c.Args().First()
		fileKey := strings.ToLower(config.View)
		config.assignDefaults()
		config.interpretLogLevel()
		if ViewMap[fileKey] != nil {
			logger.Info(ViewMap[fileKey].Data)
			return nil
		}
		return fmt.Errorf("the requested file, '%s', does not seem to exist", config.View)
	}
}
