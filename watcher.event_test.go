package main

import (
	"testing"

	"github.com/fsnotify/fsnotify"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WatcherEventTestSuite struct {
	suite.Suite
	absoluteFilePath string
	fileName         string
	fileExtension    string
}

func TestWatcherEvent(t *testing.T) {
	suite.Run(t, new(WatcherEventTestSuite))
}

func (s *WatcherEventTestSuite) SetupTest() {
	s.absoluteFilePath = "/absolute/path/to/file.ext"
	s.fileName = "file.ext"
	s.fileExtension = ".ext"
}

func (s *WatcherEventTestSuite) TestEventType_Create() {
	e := WatcherEvent(fsnotify.Event{
		Op: fsnotify.Create,
	})
	assert.Equal(s.T(), "+", e.EventType())
}

func (s *WatcherEventTestSuite) TestEventType_Write() {
	e := WatcherEvent(fsnotify.Event{
		Op: fsnotify.Write,
	})
	assert.Equal(s.T(), ">", e.EventType())
}

func (s *WatcherEventTestSuite) TestEventType_Rename() {
	e := WatcherEvent(fsnotify.Event{
		Op: fsnotify.Rename,
	})
	assert.Equal(s.T(), "/", e.EventType())
}

func (s *WatcherEventTestSuite) TestEventType_Remove() {
	e := WatcherEvent(fsnotify.Event{
		Op: fsnotify.Remove,
	})
	assert.Equal(s.T(), "-", e.EventType())
}

func (s *WatcherEventTestSuite) TestEventType_Chmod() {
	e := WatcherEvent(fsnotify.Event{
		Op: fsnotify.Chmod,
	})
	assert.Equal(s.T(), "%", e.EventType())
}

func (s *WatcherEventTestSuite) TestEventTypeString_Create() {
	e := WatcherEvent(fsnotify.Event{
		Op: fsnotify.Create,
	})
	assert.Equal(s.T(), "create", e.EventTypeString())
}

func (s *WatcherEventTestSuite) TestEventTypeString_Write() {
	e := WatcherEvent(fsnotify.Event{
		Op: fsnotify.Write,
	})
	assert.Equal(s.T(), "write", e.EventTypeString())
}

func (s *WatcherEventTestSuite) TestEventTypeString_Rename() {
	e := WatcherEvent(fsnotify.Event{
		Op: fsnotify.Rename,
	})
	assert.Equal(s.T(), "rename", e.EventTypeString())
}

func (s *WatcherEventTestSuite) TestEventTypeString_Remove() {
	e := WatcherEvent(fsnotify.Event{
		Op: fsnotify.Remove,
	})
	assert.Equal(s.T(), "remove", e.EventTypeString())
}

func (s *WatcherEventTestSuite) TestEventTypeString_Chmod() {
	e := WatcherEvent(fsnotify.Event{
		Op: fsnotify.Chmod,
	})
	assert.Equal(s.T(), "perms", e.EventTypeString())
}

func (s *WatcherEventTestSuite) TestFilePath() {
	e := WatcherEvent(fsnotify.Event{
		Op:   fsnotify.Chmod,
		Name: s.absoluteFilePath,
	})
	assert.Equal(s.T(), s.absoluteFilePath, e.FilePath())
}

func (s *WatcherEventTestSuite) TestFileName() {
	e := WatcherEvent(fsnotify.Event{
		Op:   fsnotify.Chmod,
		Name: s.absoluteFilePath,
	})
	assert.Equal(s.T(), s.fileName, e.FileName())
}

func (s *WatcherEventTestSuite) TestIsAnyOf() {
	e := WatcherEvent(fsnotify.Event{
		Op:   fsnotify.Chmod,
		Name: s.absoluteFilePath,
	})
	assert.True(s.T(), e.IsAnyOf([]string{"ext"}))
	assert.True(s.T(), e.IsAnyOf([]string{".ext"}))
}

func (s *WatcherEventTestSuite) TestFileType() {
	e := WatcherEvent(fsnotify.Event{
		Op:   fsnotify.Chmod,
		Name: s.absoluteFilePath,
	})
	assert.Equal(s.T(), s.fileExtension, e.FileType())
}
