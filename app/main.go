//go:generate go run data/generate.go
package main

import (
	"fmt"
	"strings"
	"sync"

	shellquote "github.com/kballard/go-shellquote"
)

var Version string
var Commit string

func main() {
	logger := InitLogger(&LoggerConfig{
		Name:   "main",
		Format: "production",
	})
	config := InitConfig()
	if config.RunModWatch {
		logUniversalConfigurations(logger, config)
		logWatchModeConfigurations(logger, config)
		watcher := InitWatcher(&WatcherConfig{
			FileExtensions: config.FileExtensions,
			IgnoredNames:   config.IgnoredNames,
			RefreshRate:    config.Rate,
		})
		watcher.RecursivelyWatch(config.WatchDirectory)
		var pipeline []*ExecutionGroup
		for _, execGroup := range config.ExecGroups {
			executionGroup := &ExecutionGroup{}
			var executionCommands []*Command
			commands := strings.Split(execGroup, config.CommandsDelimiter)
			for _, command := range commands {
				if sections, err := shellquote.Split(command); err != nil {
					panic(err)
				} else {
					executionCommands = append(executionCommands, &Command{
						application: sections[0],
						arguments:   sections[1:],
					})
				}
			}
			executionGroup.commands = executionCommands
			pipeline = append(pipeline, executionGroup)
		}
		runner := InitRunner(&RunnerConfig{
			pipeline: pipeline,
		})

		var wg sync.WaitGroup
		watcher.BeginWatch(&wg, func(event *WatcherEvent) bool {
			logger.Info(event)
			runner.Trigger()
			return true
		})

		logger.Info("started watcher")
		wg.Wait()
	} else if config.RunView {
		switch strings.ToLower(config.View) {
		case "dockerfile":
			logger.Info("previewing contents of Dockerfile")
			fmt.Println(DataDockerfile)
			logger.Info("end of preview for contents of Dockerfile")
		case "makefile":
			logger.Info("previewing contents of Makefile")
			fmt.Println(DataMakefile)
			logger.Info("end of preview for contents of Makefile")
		case ".dockerignore":
			logger.Info("previewing contents of .dockerignore")
			fmt.Println(DataDotDockerignore)
			logger.Info("end of preview for contents of .dockerignore")
		case ".gitignore":
			logger.Info("previewing contents of .gitignore")
			fmt.Println(DataDotGitignore)
			logger.Info("end of preview for contents of .gitignore")
		}
	} else if config.RunVersion {
		fmt.Printf("godev %s-%s\n", Version, Commit)
	}
	logger.Info("bye")
}

func logUniversalConfigurations(logger *Logger, config *Config) {
	logger.Infof("flag - init       : %v", config.RunInit)
	logger.Infof("flag - test       : %v", config.RunTest)
	logger.Infof("flag - view       : %v", config.RunView)
	logger.Infof("flag - watch      : %v", config.RunModWatch)
	logger.Infof("watch directory   : %s", config.WatchDirectory)
	logger.Infof("build output      : %s", config.BuildOutput)
}

func logWatchModeConfigurations(logger *Logger, config *Config) {
	logger.Infof("file extensions   : %v", config.FileExtensions)
	logger.Infof("ignored names     : %v", config.IgnoredNames)
	logger.Infof("refresh interval  : %v", config.Rate)
	logger.Infof("execution delim   : %s", config.CommandsDelimiter)
	logger.Info("execution groups as follows...")
	for egIndex, execGroup := range config.ExecGroups {
		logger.Infof("  %v) %s", egIndex+1, execGroup)
		commands := strings.Split(execGroup, config.CommandsDelimiter)
		for cIndex, command := range commands {
			sections, err := shellquote.Split(command)
			if err != nil {
				panic(err)
			}
			app := sections[0]
			args := sections[1:]
			logger.Infof("    %v > %s %v", cIndex+1, app, args)
		}
	}
}
