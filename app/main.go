package main

import (
	"strings"
	"sync"

	"github.com/kballard/go-shellquote"
)

func main() {
	logger := InitLogger(&LoggerConfig{
		Name:   "main",
		Format: "production",
	})
	config := InitConfig()
	logUniversalConfigurations(logger, config)
	if config.RunModWatch {
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
	}
	logger.Info("bye")
}

func logUniversalConfigurations(logger *Logger, config *Config) {
	logger.Infof("flag - init       : %v", config.RunInit)
	logger.Infof("flag - test       : %v", config.RunTest)
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
