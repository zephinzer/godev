package main

import (
	"bytes"
	"fmt"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MainTestSuite struct {
	suite.Suite
	godev *GoDev
	logs  bytes.Buffer
}

func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}

func (s *MainTestSuite) SetupTest() {
	s.godev = InitGoDev(&Config{
		CommandsDelimiter: ",",
		ExecGroups: []string{
			"echo 'a b' c,echo 'd e',echo f",
			"echo 1,echo 2 3",
			"echo ''",
		},
		LogLevel:      "trace",
		WorkDirectory: "/work/directory",
	})
	s.godev.logger.SetOutput(&s.logs)
}

func (s *MainTestSuite) Test_createPipeline_separatesCommandsCorrectly() {
	t := s.T()
	pipeline := s.godev.createPipeline()
	assert.Len(t, pipeline[0].commands, 3)
	assert.Len(t, pipeline[1].commands, 2)
	assert.Len(t, pipeline[2].commands, 1)
}

func (s *MainTestSuite) Test_createPipeline_separatesCommandArgsCorrectly() {
	t := s.T()
	pipeline := s.godev.createPipeline()
	// echo 'a b' c
	assert.Len(t, pipeline[0].commands[0].config.Arguments, 2)
	assert.Equal(t, "a b", pipeline[0].commands[0].config.Arguments[0])
	assert.Equal(t, "c", pipeline[0].commands[0].config.Arguments[1])
	// echo 'd e'
	assert.Len(t, pipeline[0].commands[1].config.Arguments, 1)
	assert.Equal(t, "d e", pipeline[0].commands[1].config.Arguments[0])
	// echo f
	assert.Len(t, pipeline[0].commands[2].config.Arguments, 1)
	assert.Equal(t, "f", pipeline[0].commands[2].config.Arguments[0])
	// echo 1
	assert.Len(t, pipeline[1].commands[0].config.Arguments, 1)
	assert.Equal(t, "1", pipeline[1].commands[0].config.Arguments[0])
	// echo 2 3
	assert.Len(t, pipeline[1].commands[1].config.Arguments, 2)
	assert.Equal(t, "2", pipeline[1].commands[1].config.Arguments[0])
	assert.Equal(t, "3", pipeline[1].commands[1].config.Arguments[1])
	// echo ''
	assert.Len(t, pipeline[2].commands[0].config.Arguments, 1)
	assert.Equal(t, "", pipeline[2].commands[0].config.Arguments[0])
}

func (s *MainTestSuite) Test_initialiseInitialisers() {
	initialisers := s.godev.initialiseInitialisers()
	assert.Len(s.T(), initialisers, 7)
	var keys []string
	for _, initialiser := range initialisers {
		keys = append(keys, initialiser.GetKey())
	}
	assert.Contains(s.T(), keys, ".git")
	assert.Contains(s.T(), keys, ".gitignore")
	assert.Contains(s.T(), keys, ".dockerignore")
	assert.Contains(s.T(), keys, "dockerfile")
	assert.Contains(s.T(), keys, "makefile")
	assert.Contains(s.T(), keys, "main.go")
	assert.Contains(s.T(), keys, "go.mod")
}

func (s *MainTestSuite) Test_initialiseRunner() {
	assert.Nil(s.T(), s.godev.runner)
	s.godev.initialiseRunner()
	assert.NotNil(s.T(), s.godev.runner)
}

func (s *MainTestSuite) Test_initialiseWatcher() {
	s.godev.config.FileExtensions = []string{"a", "b", "c"}
	s.godev.config.IgnoredNames = []string{".cache", ".git", "bin", "vendor"}
	s.godev.config.Rate = time.Second * 2
	s.godev.config.WatchDirectory = getCurrentWorkingDirectory()
	assert.Nil(s.T(), s.godev.watcher)
	s.godev.initialiseWatcher()
	assert.NotNil(s.T(), s.godev.watcher)
}

func (s *MainTestSuite) Test_initialiseWatcher_withInvalidWatchDirectory() {
	defer func() {
		r := recover()
		assert.Contains(
			s.T(),
			string(r.(string)),
			"/does/and/should/not/exist' does not exist",
		)
	}()
	s.godev.config.FileExtensions = []string{"a", "b", "c"}
	s.godev.config.IgnoredNames = []string{"d", "e", "f"}
	s.godev.config.Rate = time.Second * 2
	s.godev.config.WatchDirectory =
		path.Join(
			getCurrentWorkingDirectory(),
			"/does/and/should/not/exist",
		)
	assert.Nil(s.T(), s.godev.watcher)
	s.godev.initialiseWatcher()
}

func (s *MainTestSuite) Test_logUniversalConfiguration() {
	s.godev.logUniversalConfigurations()
	logs := s.logs.String()
	assert.Contains(s.T(), logs, "flag - init")
	assert.Contains(s.T(), logs, "flag - test")
	assert.Contains(s.T(), logs, "flag - view")
	assert.Contains(s.T(), logs, "watch directory")
	assert.Contains(s.T(), logs, "work directory")
	assert.Contains(s.T(), logs, "build output")
}

func (s *MainTestSuite) Test_viewFile_thatExists() {
	s.godev.config.View = "main.go"
	s.godev.viewFile()
	assert.Contains(s.T(), s.logs.String(), "previewing contents of main.go")
}

func (s *MainTestSuite) Test_viewFile_thatDoesntExist() {
	defer func() {
		r := recover()
		err := fmt.Sprintf("%s", r)
		assert.Contains(s.T(), err, "file 'nonexistent.file' does not seem to exist")
	}()
	s.godev.config.View = "nonexistent.file"
	s.godev.viewFile()
}
