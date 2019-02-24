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
	if config.RunModWatch {
		watcher = InitWatcher(&WatcherConfig{
			FileExtensions: config.FileExtensions,
			IgnoredNames:   config.IgnoredNames,
			RefreshRate:    config.Rate,
		})
		logger.Infof("file extensions  : %v", watcher.config.FileExtensions)
		logger.Infof("ignored names    : %v", watcher.config.IgnoredNames)
		logger.Infof("refresh interval : %v", watcher.config.RefreshRate)
	}
	watcher.RecursivelyWatch(config.WatchDirectory)
	var wg sync.WaitGroup
	// tick := time.Tick(5 * time.Second)
	// go func() {
	// 	for {
	// 		select {
	// 		case <-tick:
	// 			watcher.EndWatch()
	// 		}
	// 	}
	// }()
	watcher.BeginWatch(&wg, func(event *WatcherEvent) bool {
		logger.Info(event)
		return true
	})
	logger.Info("started watcher")
	wg.Wait()
	logger.Info("bye")
}
