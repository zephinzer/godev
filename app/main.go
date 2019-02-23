package main

import (
	"fmt"
	"os"
)

func main() {
	watcher := InitWatcher(&WatcherConfig{
		FileExtensions: []string{"go"},
		IgnoredNames:   []string{"vendor", ".cache"},
	})
	if cwd, err := os.Getwd(); err != nil {
		panic(err)
	} else {
		watcher.RecursivelyWatch(cwd)
	}
	config := InitConfig()
	fmt.Println(config)
}
