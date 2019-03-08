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
	pipeline := s.godev.createPipeline()
	assert.Len(s.T(), pipeline[0].commands, 3)
	assert.Len(s.T(), pipeline[1].commands, 2)
	assert.Len(s.T(), pipeline[2].commands, 1)
}

func (s *MainTestSuite) Test_createPipeline_separatesCommandArgsCorrectly() {
	pipeline := s.godev.createPipeline()
	// echo 'a b' c
	assert.Len(s.T(), pipeline[0].commands[0].config.Arguments, 2)
	assert.Equal(s.T(), "a b", pipeline[0].commands[0].config.Arguments[0])
	assert.Equal(s.T(), "c", pipeline[0].commands[0].config.Arguments[1])
	// echo 'd e'
	assert.Len(s.T(), pipeline[0].commands[1].config.Arguments, 1)
	assert.Equal(s.T(), "d e", pipeline[0].commands[1].config.Arguments[0])
	// echo f
	assert.Len(s.T(), pipeline[0].commands[2].config.Arguments, 1)
	assert.Equal(s.T(), "f", pipeline[0].commands[2].config.Arguments[0])
	// echo 1
	assert.Len(s.T(), pipeline[1].commands[0].config.Arguments, 1)
	assert.Equal(s.T(), "1", pipeline[1].commands[0].config.Arguments[0])
	// echo 2 3
	assert.Len(s.T(), pipeline[1].commands[1].config.Arguments, 2)
	assert.Equal(s.T(), "2", pipeline[1].commands[1].config.Arguments[0])
	assert.Equal(s.T(), "3", pipeline[1].commands[1].config.Arguments[1])
	// echo ''
	assert.Len(s.T(), pipeline[2].commands[0].config.Arguments, 1)
	assert.Equal(s.T(), "", pipeline[2].commands[0].config.Arguments[0])
}

func (s *MainTestSuite) Test_initialiseRunner() {
	assert.Nil(s.T(), s.godev.runner)
	s.godev.initialiseRunner()
	assert.NotNil(s.T(), s.godev.runner)
}

func (s *MainTestSuite) Test_initialiseWatcher() {
	s.godev.config.FileExtensions = []string{"a", "b", "c"}
	s.godev.config.IgnoredNames = []string{"d", "e", "f"}
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
