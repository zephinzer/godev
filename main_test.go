package main

import (
	"bytes"
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
		CommandArguments:  []string{"test", "arg"},
		CommandsDelimiter: ",",
		ExecGroups: []string{
			"echo 'a b' c,echo 'd e',echo f",
			"echo 1,echo 2 3",
			"echo ''",
		},
		EnvVars:       []string{"A=1", "B=2"},
		LogLevel:      "trace",
		WorkDirectory: "/work/directory",
	})
	s.godev.logger.SetOutput(&s.logs)
}

func (s *MainTestSuite) Test_createPipeline_assignsEnvVarsCorrectly() {
	t := s.T()
	pipeline := s.godev.createPipeline()
	for _, executionGroup := range pipeline {
		for _, command := range executionGroup.commands {
			assert.Len(t, command.config.Environment, 2)
		}
	}
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
	assert.Len(t, pipeline[2].commands[0].config.Arguments, 3)
	assert.Equal(t, "", pipeline[2].commands[0].config.Arguments[0])
	assert.Equal(t, "test", pipeline[2].commands[0].config.Arguments[1])
	assert.Equal(t, "arg", pipeline[2].commands[0].config.Arguments[2])
}

func (s *MainTestSuite) Test_eventHandler() {
	t := s.T()
	// set exec groups to none so that no pipeline triggers
	s.godev.config.ExecGroups = []string{}
	s.godev.initialiseRunner()
	s.godev.eventHandler(&[]WatcherEvent{
		WatcherEvent{Op: 1},
		WatcherEvent{Op: 2},
		WatcherEvent{Op: 4},
		WatcherEvent{Op: 8},
		WatcherEvent{Op: 16},
	})
	logs := s.logs.String()
	assert.Contains(t, logs, "CREATE")
	assert.Contains(t, logs, "WRITE")
	assert.Contains(t, logs, "REMOVE")
	assert.Contains(t, logs, "RENAME")
	assert.Contains(t, logs, "CHMOD")
}

func (s *MainTestSuite) Test_initialiseInitialisers() {
	t := s.T()
	initialisers := s.godev.initialiseInitialisers()
	assert.Len(t, initialisers, 7)
	var keys []string
	for _, initialiser := range initialisers {
		keys = append(keys, initialiser.GetKey())
	}
	assert.Contains(t, keys, ".git")
	assert.Contains(t, keys, ".gitignore")
	assert.Contains(t, keys, ".dockerignore")
	assert.Contains(t, keys, "dockerfile")
	assert.Contains(t, keys, "makefile")
	assert.Contains(t, keys, "main.go")
	assert.Contains(t, keys, "go.mod")
}

func (s *MainTestSuite) Test_initialiseRunner() {
	t := s.T()
	assert.Nil(t, s.godev.runner)
	s.godev.initialiseRunner()
	assert.NotNil(t, s.godev.runner)
}

func (s *MainTestSuite) Test_initialiseWatcher() {
	t := s.T()
	s.godev.config.FileExtensions = []string{"a", "b", "c"}
	s.godev.config.IgnoredNames = []string{".cache", ".git", "bin", "vendor"}
	s.godev.config.Rate = time.Second * 2
	s.godev.config.WatchDirectory = getCurrentWorkingDirectory()
	assert.Nil(t, s.godev.watcher)
	s.godev.initialiseWatcher()
	assert.NotNil(t, s.godev.watcher)
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
	t := s.T()
	s.godev.logUniversalConfigurations()
	logs := s.logs.String()
	assert.Contains(t, logs, "flag - init")
	assert.Contains(t, logs, "flag - test")
	assert.Contains(t, logs, "flag - view")
	assert.Contains(t, logs, "watch directory")
	assert.Contains(t, logs, "work directory")
	assert.Contains(t, logs, "build output")
}

func (s *MainTestSuite) Test_logWatchModeConfigurations() {
	t := s.T()
	s.godev.logWatchModeConfigurations()
	logs := s.logs.String()
	assert.Contains(t, logs, "environment")
	assert.Contains(t, logs, "file extensions")
	assert.Contains(t, logs, "ignored names")
	assert.Contains(t, logs, "refresh interval")
	assert.Contains(t, logs, "execution delim")
	assert.Contains(t, logs, "execution groups")
	assert.Contains(t, logs, "1) echo 'a b' c,echo 'd e',echo f")
	assert.Contains(t, logs, "2) echo 1,echo 2 3")
	assert.Contains(t, logs, "3) echo ''")
	assert.Contains(t, logs, "test arg")
}
