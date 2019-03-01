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

var watcherEventTypeString = []string{
	"",
	// 1
	"create",
	// 2
	"write", "",
	// 4
	"remove", "", "", "",
	// 8
	"rename", "", "", "", "", "", "", "",
	//  16
	"perms",
}

// WatcherEvent provides some function candy for working with
// fsnotify more easily
type WatcherEvent fsnotify.Event

// EventType returns a symbol denoting the type of operation recorded
func (e *WatcherEvent) EventType() string {
	return watcherEventType[e.Op]
}

// EventTypeString returns a string denoting the type of operation recorded
func (e *WatcherEvent) EventTypeString() string {
	return watcherEventTypeString[e.Op]
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
		return WatcherFileTypeDeleted
	default:
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
