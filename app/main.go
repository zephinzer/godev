package main

import (
	"sync"
)

func main() {
	logger := InitLogger(&LoggerConfig{
		Name:   "main",
		Format: "production",
	})
	config := InitConfig()
	var watcher *Watcher
	var runner *Runner
	logger.Infof("flag[init]       : %v", config.RunInit)
	logger.Infof("flag[test]       : %v", config.RunTest)
	logger.Infof("flag[watch]      : %v", config.RunModWatch)
	if config.RunModWatch {
		watcher = InitWatcher(&WatcherConfig{
			FileExtensions: config.FileExtensions,
			IgnoredNames:   config.IgnoredNames,
			RefreshRate:    config.Rate,
		})
		logger.Infof("file extensions  : %v", watcher.config.FileExtensions)
		logger.Infof("ignored names    : %v", watcher.config.IgnoredNames)
		logger.Infof("refresh interval : %v", watcher.config.RefreshRate)
		watcher.RecursivelyWatch(config.WatchDirectory)

		runner = InitRunner(&RunnerConfig{
			pipeline: []ExecutionGroup{
				ExecutionGroup{
					commands: []*Command{
						&Command{
							application: "./data/test-run/sleeper.sh",
							arguments:   []string{"1"},
						},
						&Command{
							application: "ls",
							arguments:   []string{"-a"},
						},
					},
				},
			},
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
