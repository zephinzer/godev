package main

import (
	"fmt"
	"os"
	"path"
	"strings"

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
	// WatcherFileTypeDir indicates a directory
	WatcherFileTypeDir = "dir"
	// WatcherFileTypeErrored indicates an error
	WatcherFileTypeErrored = "err"
	// WatcherFileTypeDeleted indicates a deleted item
	WatcherFileTypeDeleted = "rm"
)

var watcherEventType = []string{
	"",
	// 00001
	WatcherEventCreate,
	// 00010
	WatcherEventWrite, "",
	// 00100
	WatcherEventRemove, "", "", "",
	// 01000
	WatcherEventRename, "", "", "", "", "", "", "",
	// 10000
	WatcherEventPermission,
}

// WatcherEvent provides some function candy for working with
// fsnotify more easily
type WatcherEvent fsnotify.Event

// EventType returns a symbol denoting the type of operation recorded
func (e *WatcherEvent) EventType() string {
	eventType := ""
	watcherEvents := []fsnotify.Op{1, 2, 4, 8, 16}
	for _, event := range watcherEvents {
		if e.Op|event == event {
			eventType += watcherEventType[event]
		}
	}
	return eventType
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
	if e.Op|fsnotify.Remove == fsnotify.Remove {
		return WatcherFileTypeDeleted
	} else {
		fileType := path.Ext(e.Name)
		if len(fileType) == 0 {
			fileInfo, err := os.Lstat(e.Name)
			if err != nil {
				return WatcherFileTypeErrored
			} else if fileInfo.IsDir() {
				return WatcherFileTypeDir
			} else {
				return path.Base(e.Name)
			}
		} else {
			return fileType
		}
	}
}

// IsAnyOf verifies that the file extension matches :theseTypes
func (e *WatcherEvent) IsAnyOf(theseTypes []string) bool {
	for _, fileExtension := range theseTypes {
		if strings.TrimLeft(e.FileType(), ".") == strings.TrimLeft(fileExtension, ".") {
			return true
		}
	}
	return false
}

func (e *WatcherEvent) String() string {
	return fmt.Sprintf(
		"[%s] %s at '%s'",
		e.EventType(),
		e.FileType(),
		e.FilePath(),
	)
}
