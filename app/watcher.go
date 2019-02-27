package main

import (
	"fmt"
	"io/ioutil"
	_ "log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// WatcherConfig is for configuring Watcher
type WatcherConfig struct {
	FileExtensions []string
	IgnoredNames   []string
	RefreshRate    time.Duration
	LogLevel       LogLevel
}

// InitWatcher returns a workable Watcher instance
func InitWatcher(config *WatcherConfig) *Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	fw := &Watcher{
		config:  config,
		logger:  InitLogger(&LoggerConfig{Name: "watcher", Format: "production", Level: config.LogLevel}),
		watcher: watcher,
	}
	return fw
}

// Watcher is a component for handling file system changes
type Watcher struct {
	config         *WatcherConfig
	logger         *Logger
	watcher        *fsnotify.Watcher
	events         []WatcherEvent
	watchMutex     chan bool
	intervalTicker <-chan time.Time
}

// Close closes the watcher, use for graceful shutdowns
func (fw *Watcher) Close() {
	if fw.watcher == nil {
		panic("watcher was not initialised")
	}
	fw.watcher.Close()
}

// WatcherEventHandler defines the callback for BeginWatch() to use
type WatcherEventHandler func(*[]WatcherEvent) bool

// BeginWatch starts the file system watching in blocking mode
func (fw *Watcher) BeginWatch(waitGroup *sync.WaitGroup, handler WatcherEventHandler) {
	fw.logger.Trace("initialising file system watch")
	fw.watchMutex = make(chan bool)
	fw.intervalTicker = time.After(fw.config.RefreshRate)
	waitGroup.Add(1)
	go fw.watchRoutine(
		fw.intervalTicker,
		fw.watchMutex,
		handler,
		waitGroup.Done,
	)
}

// EndWatch sends the signal to main goroutine to end itself
func (fw *Watcher) EndWatch() {
	fw.logger.Trace("received signal to terminate file system watch")
	fw.watchMutex <- true
}

func (fw *Watcher) watchRoutine(tick <-chan time.Time, stop chan bool, handler WatcherEventHandler, onDone func()) {
	for {
		select {
		case <-tick:
			if len(fw.events) > 0 {
				fw.logger.Tracef("processing %v raw events...", len(fw.events))
				dedupedEvents := fw.getDedupedEvents()
				handler(&dedupedEvents)
				fw.logger.Tracef("processed %v event(s)", len(dedupedEvents))
				fw.events = make([]WatcherEvent, 0)
			}
		case event := <-fw.watcher.Events:
			eventToAdd := WatcherEvent(event)
			if eventToAdd.IsAnyOf(fw.config.FileExtensions) {
				fw.events = append(fw.events, eventToAdd)
				tick = time.After(2 * time.Second)
			}
		case shouldWeStop := <-stop:
			fw.logger.Tracef("received signal to terminate watch routine: %v", shouldWeStop)
			fw.watchMutex = make(chan bool)
			onDone()
			if shouldWeStop {
				break
			}
		}
	}
}

// RecursivelyWatch is so we can watch all sub directories of a directory
func (fw *Watcher) RecursivelyWatch(directoryPath string) {
	fw.assertDirectoryIntegrity(directoryPath)
	allSubDirectories := fw.recursivelyGetDirectories(directoryPath)
	fw.Watch(directoryPath)
	for _, directory := range allSubDirectories {
		fw.Watch(directory)
	}
}

// Watch is here for watching a single directory
func (fw *Watcher) Watch(directoryPath string) {
	fw.assertDirectoryIntegrity(directoryPath)
	fw.watcher.Add(directoryPath)
	fw.logger.Tracef("registered '%s'", directoryPath)
}

// assertDirectoryIntegrity panicks if the :directoryPath does not exist/is not a directory
func (fw *Watcher) assertDirectoryIntegrity(directoryPath string) {
	if !fw.pathExists(directoryPath) {
		panic(fmt.Sprintf("provided path '%s' does not exist", directoryPath))
	} else if !fw.pathIsDirectory(directoryPath) {
		panic(fmt.Sprintf("provided path '%s' is not a directory", directoryPath))
	}
}

// getDedupedEvents processes the events so that we don't respond to duplicate items
func (fw *Watcher) getDedupedEvents() []WatcherEvent {
	var eventsProcessed []string
	var eventsToProcess []WatcherEvent
	for _, event := range fw.events {
		proceed := true
		for _, processedEvent := range eventsProcessed {
			if processedEvent == event.String() {
				proceed = false
			}
		}
		if proceed {
			eventsProcessed = append(eventsProcessed, event.String())
			eventsToProcess = append(eventsToProcess, event)
		}
	}
	return eventsToProcess
}

// isIgnoredName checks whether the name was faulty
func (fw *Watcher) isIgnoredName(name string) bool {
	ignore := false
	if fw.config == nil {
		return ignore
	}
	for _, ignoredName := range fw.config.IgnoredNames {
		if ignoredName == name {
			ignore = true
		}
	}
	return ignore
}

// pathIsDirectory is for argument verification
func (fw *Watcher) pathExists(absolutePath string) bool {
	if _, err := os.Lstat(absolutePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// pathIsDirectory is for argument verification
func (fw *Watcher) pathIsDirectory(absolutePath string) bool {
	if fileInfo, err := os.Lstat(absolutePath); err != nil {
		panic(err)
	} else {
		return fileInfo.IsDir()
	}
}

// recursivelyGetDirectories is here to retrieve a list of all sub-directories from :directoryPath
func (fw *Watcher) recursivelyGetDirectories(directoryPath string) []string {
	fw.assertDirectoryIntegrity(directoryPath)
	directoryListing, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		panic(err)
	}
	var listings []string
	for _, listing := range directoryListing {
		listingFullPath := path.Join(directoryPath, listing.Name())
		if !fw.isIgnoredName(listing.Name()) && listing.IsDir() {
			listings = append(listings, listingFullPath)
			listings = append(listings, fw.recursivelyGetDirectories(listingFullPath)...)
		}
	}
	return listings
}
