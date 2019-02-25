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
	logger.Infof("flag - init      : %v", config.RunInit)
	logger.Infof("flag - test      : %v", config.RunTest)
	logger.Infof("flag - watch     : %v", config.RunModWatch)
	if config.RunModWatch {
		watcher := InitWatcher(&WatcherConfig{
			FileExtensions: config.FileExtensions,
			IgnoredNames:   config.IgnoredNames,
			RefreshRate:    config.Rate,
		})
		watcher.RecursivelyWatch(config.WatchDirectory)
		runner := InitRunner(&RunnerConfig{
			pipeline: []*ExecutionGroup{
				&ExecutionGroup{
					commands: []*Command{
						&Command{
							application: "./data/test-run/sleeper.sh",
							arguments:   []string{"1"},
						},
						&Command{
							application: "echo",
							arguments:   []string{"a"},
						},
					},
				},
				&ExecutionGroup{
					commands: []*Command{
						&Command{
							application: "./data/test-run/sleeper.sh",
							arguments:   []string{"1"},
						},
						&Command{
							application: "./data/test-run/sleeper.sh",
							arguments:   []string{"1"},
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
