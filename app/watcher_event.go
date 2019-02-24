package main

import (
	"fmt"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
)

const (
	// WatcherEventCreate denotes a file/dir creation
	WatcherEventCreate = "+"
	// WatcherEventWrite denotes a file write
	WatcherEventWrite = ">"
	// WatcherEventRemove denotes removal of a file/dir
	WatcherEventRemove = "-"
	// WatcherEventRename denotes renaming of a file/dir
	WatcherEventRename = "/"
	// WatcherEventPermission denotes chmodding of a file/dir
	WatcherEventPermission = "%"
)

var watcherEventType = []string{
	"",
	// 1
	WatcherEventCreate,
	// 2
	WatcherEventWrite, "",
	// 4
	WatcherEventRemove, "", "", "",
	// 8
	WatcherEventRename, "", "", "", "", "", "", "",
	//  16
	WatcherEventPermission,
}

// WatcherEvent provides some function candy for working with
// fsnotify more easily
type WatcherEvent fsnotify.Event

// EventType returns a symbol denoting the type of operation recorded
func (e *WatcherEvent) EventType() string {
	return watcherEventType[e.Op]
}

// FilePath returns the absolute path of the file/dir
func (e *WatcherEvent) FilePath() string {
	return e.Name
}

// FileName returns the file/dir name
func (e *WatcherEvent) FileName() string {
	return path.Base(e.Name)
}

// FileType returns the extension of the file if its a file,
// "dir" if its a dir, or "errored" if an error occurred
func (e *WatcherEvent) FileType() string {
	switch e.EventType() {
	case WatcherEventRemove:
		return "deleted"
	default:
		fileType := path.Ext(e.Name)
		if len(fileType) == 0 {
			fileInfo, err := os.Lstat(e.Name)
			if err != nil {
				fileType = "errored"
			} else if fileInfo.IsDir() {
				fileType = "dir"
			} else {
				fileType = path.Base(e.Name)
			}
		}
		return fileType
	}
}

func (e *WatcherEvent) String() string {
	return fmt.Sprintf("%s: [%s] at '%s'", e.EventType(), e.FileType(), e.FilePath())
}
