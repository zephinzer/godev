package main

import (
	"fmt"
	"os"
	"path"
	"testing"

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

func (s *WatcherTestSuite) TestWatch() {
	w := InitWatcher(&WatcherConfig{})
	cwd := s.currentDirectory
	testDirectoryPath := path.Join(cwd, "/data/test-watch")
	w.Watch(testDirectoryPath)
}

func (s *WatcherTestSuite) Test_assertDirectoryIntegrityPass() {
	defer func() {
		err := recover()
		assert.NotNil(s.T(), err, "expected an error but none was panicked")
	}()
	w := &Watcher{}
	cwd := s.currentDirectory
	w.assertDirectoryIntegrity(path.Join(cwd, "/non/existent"))
}

func (s *WatcherTestSuite) Test_assertDirectoryIntegrityFail() {
	defer func() {
		err := recover()
		assert.Nilf(s.T(), err, "expected no errors but '%s' was panicked", err)
	}()
	w := &Watcher{}
	cwd := s.currentDirectory
	w.assertDirectoryIntegrity(path.Join(cwd))
}

func (s *WatcherTestSuite) Test_isIgnoredName() {
	ignoredName := "ignored"
	notIgnoredNames := []string{
		fmt.Sprintf(" %s", ignoredName),
		fmt.Sprintf("%s ", ignoredName),
		fmt.Sprintf(" %s ", ignoredName),
	}
	w := &Watcher{
		config: &WatcherConfig{
			IgnoredNames: []string{ignoredName},
		},
	}
	assert.Truef(s.T(), w.isIgnoredName(ignoredName), "expected '%s' to be ignored but it wasn't", ignoredName)
	for _, nameToWatch := range notIgnoredNames {
		assert.Falsef(s.T(), w.isIgnoredName(nameToWatch), "expected '%s' to not be ignored but it was", nameToWatch)
	}
}

func (s *WatcherTestSuite) Test_pathIsDirectory() {
	w := &Watcher{}
	cwd := s.currentDirectory
	assert.Truef(s.T(), w.pathIsDirectory(cwd), "expected '%s' (current working directory) to be a directory but it was not", cwd)
}

func (s *WatcherTestSuite) Test_pathExists() {
	w := &Watcher{}
	cwd := s.currentDirectory
	assert.Truef(s.T(), w.pathExists(cwd), "expected '%s' (current working directory) to exist but it did not", cwd)
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
