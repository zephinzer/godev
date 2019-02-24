package main

import (
	"sync"
)

func main() {
	logger := InitLogger(&LoggerConfig{
		Name: "main",
	})
	logger.Info("hello")
	config := InitConfig()
	var watcher *Watcher
	if config.RunModWatch {
		watcher = InitWatcher(&WatcherConfig{
			FileExtensions: []string{"go"},
			IgnoredNames:   []string{"vendor", ".cache"},
			RefreshRate:    config.Rate,
		})
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
		logger.Info(event.String())
		return true
	})
	logger.Info("bye")
}
