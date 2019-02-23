package main

import (
	"fmt"
	"io/ioutil"
	_ "log"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
)

// WatcherConfig is for configuring Watcher
type WatcherConfig struct {
	FileExtensions []string
	IgnoredNames   []string
}

func InitWatcher(config *WatcherConfig) *Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	fw := &Watcher{
		config:  config,
		logger:  InitLogger(&LoggerConfig{Name: "watcher", Format: "text", Level: "trace"}),
		watcher: watcher,
	}
	return fw
}

type WatcherEvent struct {
}

// Watcher is a component for handling file system changes
type Watcher struct {
	config      *WatcherConfig
	logger      *Logger
	watcher     *fsnotify.Watcher
	eventsQueue chan fsnotify.Event
	events      chan []WatcherEvent
}

// Close closes the watcher, use for graceful shutdowns
func (fw *Watcher) Close() {
	if fw.watcher == nil {
		panic("watcher was not initialised")
	}
	fw.watcher.Close()
}

func (fw *Watcher) RecursivelyWatch(directoryPath string) {
	fw.assertDirectoryIntegrity(directoryPath)
	allSubDirectories := fw.recursivelyGetDirectories(directoryPath)
	for _, directory := range allSubDirectories {
		fw.Watch(directory)
	}
}

// Watch is here for watching a single directory
func (fw *Watcher) Watch(directoryPath string) {
	fw.assertDirectoryIntegrity(directoryPath)
	fw.logger.Tracef("registering '%s'", directoryPath)
	fw.watcher.Add(directoryPath)
}

// assertDirectoryIntegrity panicks if the :directoryPath does not exist/is not a directory
func (fw *Watcher) assertDirectoryIntegrity(directoryPath string) {
	if !fw.pathExists(directoryPath) {
		panic(fmt.Sprintf("provided path '%s' does not exist", directoryPath))
	} else if !fw.pathIsDirectory(directoryPath) {
		panic(fmt.Sprintf("provided path '%s' is not a directory", directoryPath))
	}
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

func (fw *Watcher) pathIsDirectory(absolutePath string) bool {
	if fileInfo, err := os.Lstat(absolutePath); err != nil {
		panic(err)
	} else {
		return fileInfo.IsDir()
	}
}

func (fw *Watcher) pathExists(absolutePath string) bool {
	if _, err := os.Lstat(absolutePath); os.IsNotExist(err) {
		return false
	}
	return true
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
