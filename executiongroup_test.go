package main

import (
	"bytes"
	"regexp"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ExecutionGroupTestSuite struct {
	suite.Suite
	executionGroup *ExecutionGroup
	logs           bytes.Buffer
	logger         *Logger
}

func TestExecutionGroup(t *testing.T) {
	suite.Run(t, new(ExecutionGroupTestSuite))
}

func (s *ExecutionGroupTestSuite) SetupTest() {
	logger := InitLogger(&LoggerConfig{
		Name:   "ExecutionGroupTestSuite",
		Format: "production",
		Level:  "trace",
	})
	logger.SetOutput(&s.logs)
	s.executionGroup = &ExecutionGroup{
		logger: logger,
	}
}

func (s *ExecutionGroupTestSuite) TestIsRunning() {
	t := s.T()
	s.executionGroup.commands = []*Command{
		mockCommand("echo", []string{"1"}, &s.logs),
	}
	s.executionGroup.commands[0].started = false
	assert.False(t, s.executionGroup.IsRunning())
	s.executionGroup.commands[0].started = true
	assert.True(t, s.executionGroup.IsRunning())
	s.executionGroup.commands[0].stopped = true
	assert.False(t, s.executionGroup.IsRunning())
}

func (s *ExecutionGroupTestSuite) TestRun() {
	s.executionGroup.commands = []*Command{
		mockCommand("echo", []string{"1"}, &s.logs),
		mockCommand("echo", []string{"2"}, &s.logs),
		mockCommand("echo", []string{"3"}, &s.logs),
	}
	s.executionGroup.logger.SetOutput(&s.logs)
	t := s.T()
	s.executionGroup.Run()
	assert.Contains(t, s.logs.String(), "command[echo[1]] is starting")
	assert.Contains(t, s.logs.String(), "command[echo[2]] is starting")
	assert.Contains(t, s.logs.String(), "command[echo[3]] is starting")
	assert.Regexp(t, regexp.MustCompile(`execution group\[\d\] is starting`), s.logs.String())
	assert.Regexp(t, regexp.MustCompile(`execution group\[\d\] exited`), s.logs.String())
}

func (s *ExecutionGroupTestSuite) TestTerminate() {
	t := s.T()
	s.executionGroup.commands = []*Command{
		mockCommand("echo", []string{"1"}, &s.logs),
	}
	s.executionGroup.commands[0].handleInitialisation()
	s.executionGroup.commands[0].started = true
	s.executionGroup.commands[0].stopped = false
	go func() {
		select {
		case signal := <-s.executionGroup.commands[0].signal:
			assert.Equal(t, signal, syscall.SIGINT)
		}
	}()
	s.executionGroup.Terminate()
}

func (s *ExecutionGroupTestSuite) Test_handleCommandStatus() {
	t := s.T()
	testCommand := mockCommand("echo", []string{"1"}, &s.logs)
	s.executionGroup.waitGroup.Add(1)
	s.executionGroup.handleCommandStatus(testCommand, nil)
	s.executionGroup.waitGroup.Wait()
	assert.Contains(t, s.logs.String(), "command[echo[1]] exited without error")
}
