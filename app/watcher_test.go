package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"sync"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WatcherTestSuite struct {
	suite.Suite
	currentDirectory string
}

func TestWatcher(t *testing.T) {
	suite.Run(t, new(WatcherTestSuite))
}

func (s *WatcherTestSuite) SetupTest() {
	cwd, err := os.Getwd()
	if err != nil {
		s.T().Errorf("error while retrieving current directory: %s", err)
	}
	s.currentDirectory = cwd
}

func (s *WatcherTestSuite) TestEndWatch() {
	var logBuffer bytes.Buffer
	mockLog := InitLogger(&LoggerConfig{
		Name: "test",
	})
	mockLog.SetOutput(&logBuffer)
	w := &Watcher{
		logger:     mockLog,
		watchMutex: make(chan bool, 1),
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func(tick <-chan time.Time) {
		for {
			select {
			case <-tick:
				w.EndWatch()
			case done := <-w.watchMutex:
				assert.True(s.T(), done)
				wg.Done()
			}
		}
	}(time.After(10 * time.Millisecond))
	wg.Wait()
}

func (s *WatcherTestSuite) TestRecursivelyWatch() {
	var logBuffer bytes.Buffer
	mockLog := InitLogger(&LoggerConfig{
		Name: "test",
	})
	mockLog.SetOutput(&logBuffer)
	w := InitWatcher(&WatcherConfig{})
	w.logger = mockLog
	cwd := s.currentDirectory
	testDirectoryPath := path.Join(cwd, "/data/test-recursive")
	w.RecursivelyWatch(testDirectoryPath)
	defer w.Close()
	allSubDirectories := w.recursivelyGetDirectories(testDirectoryPath)
	logs := string(logBuffer.Bytes())
	for _, subDirectory := range allSubDirectories {
		assert.Containsf(
			s.T(),
			logs,
			fmt.Sprintf("registered '%s'", subDirectory),
			"expected a line of log for registration of '%s' but none was found",
			subDirectory,
		)
	}
}

func (s *WatcherTestSuite) TestWatch() {
	t := s.T()
	var logBuffer bytes.Buffer
	mockLog := InitLogger(&LoggerConfig{
		Name: "test",
	})
	mockLog.SetOutput(&logBuffer)
	w := InitWatcher(&WatcherConfig{})
	w.logger = mockLog
	cwd := s.currentDirectory
	testDirectoryPath := path.Join(cwd, "/data/test-watch")
	w.Watch(testDirectoryPath)
	defer w.Close()
	testFilePath := path.Join(testDirectoryPath, "Watcher.TestWatch")
	createFile(t, testFilePath)
	removeFile(t, testFilePath)
	assert.Contains(
		t,
		string(logBuffer.Bytes()),
		fmt.Sprintf("registered '%s'", testDirectoryPath),
		"expected to have a line of logs indicating directory registration",
	)
}

func (s *WatcherTestSuite) Test_assertDirectoryIntegrityPass() {
	defer expectError(s.T())()
	w := &Watcher{}
	cwd := s.currentDirectory
	w.assertDirectoryIntegrity(path.Join(cwd, "/non/existent"))
}

func (s *WatcherTestSuite) Test_assertDirectoryIntegrityFail() {
	defer expectNoError(s.T())()
	w := &Watcher{}
	cwd := s.currentDirectory
	w.assertDirectoryIntegrity(path.Join(cwd))
}

func (s *WatcherTestSuite) Test_getDedupedEvents() {
	testNameDifferences := []WatcherEvent{
		WatcherEvent{Name: "/path/to/watched/file/1", Op: fsnotify.Create},
		WatcherEvent{Name: "/path/to/watched/file/1", Op: fsnotify.Create},
		WatcherEvent{Name: "/path/to/watched/file/2", Op: fsnotify.Create},
	}
	w := &Watcher{events: testNameDifferences}
	assert.Len(s.T(), w.getDedupedEvents(), 2)

	testOpDifferences := []WatcherEvent{
		WatcherEvent{Name: "/path/to/watched/file/1", Op: fsnotify.Create},
		WatcherEvent{Name: "/path/to/watched/file/1", Op: fsnotify.Create},
		WatcherEvent{Name: "/path/to/watched/file/1", Op: fsnotify.Remove},
	}
	w.events = testOpDifferences
	assert.Len(s.T(), w.getDedupedEvents(), 2)
}

func (s *WatcherTestSuite) Test_isIgnoredName() {
	ignoredName := "ignored"
	watchedNames := []string{
		fmt.Sprintf(" %s", ignoredName),
		fmt.Sprintf("%s ", ignoredName),
		fmt.Sprintf(" %s ", ignoredName),
		fmt.Sprintf("_%s", ignoredName),
		fmt.Sprintf("%s_", ignoredName),
		fmt.Sprintf("_%s_", ignoredName),
		fmt.Sprintf("a%s", ignoredName),
		fmt.Sprintf("%sa", ignoredName),
		fmt.Sprintf("a%sa", ignoredName),
		fmt.Sprintf("0%s", ignoredName),
		fmt.Sprintf("%s0", ignoredName),
		fmt.Sprintf("0%s0", ignoredName),
	}
	w := &Watcher{
		config: &WatcherConfig{
			IgnoredNames: []string{ignoredName},
		},
	}
	assert.Truef(s.T(), w.isIgnoredName(ignoredName), "expected '%s' to be ignored but it wasn't", ignoredName)
	for _, nameToWatch := range watchedNames {
		assert.Falsef(s.T(), w.isIgnoredName(nameToWatch), "expected '%s' to not be ignored but it was", nameToWatch)
	}
}

func (s *WatcherTestSuite) Test_pathExists() {
	w := &Watcher{}
	cwd := s.currentDirectory
	assert.Truef(s.T(), w.pathExists(cwd), "expected '%s' (current working directory) to exist but it did not", cwd)

	errd := path.Join(s.currentDirectory, "/____non-existent/____path")
	assert.Falsef(s.T(), w.pathExists(errd), "expected '%s' (non-existent path) to not exist but it did apparently", errd)
}

func (s *WatcherTestSuite) Test_pathIsDirectory() {
	w := &Watcher{}
	cwd := s.currentDirectory
	assert.Truef(s.T(), w.pathIsDirectory(cwd), "expected '%s' (current working directory) to be a directory but it was not", cwd)

	filePath := path.Join(s.currentDirectory, "/main.go")
	assert.Falsef(s.T(), w.pathIsDirectory(filePath), "expected '%s' to not be a directory but it was", filePath)
}

func (s *WatcherTestSuite) Test_recursivelyGetDirectories() {
	w := &Watcher{}
	cwd := s.currentDirectory
	expectedDirectories := []string{"1", "2", "2-1", "2-2", "2-2-1", "3"}
	directories := w.recursivelyGetDirectories(path.Join(cwd, "/data/test-recursive"))
	for index, directory := range directories {
		assert.Equalf(s.T(), path.Base(directory), expectedDirectories[index], "expected '%s' to be '%s", path.Base(directory), expectedDirectories[index])
	}
}
