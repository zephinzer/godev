package main

import (
	"bytes"
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
	s.executionGroup = &ExecutionGroup{
		commands: []ICommand{
			InitCommandMock("echo", []string{"1"}, &s.logs),
			InitCommandMock("echo", []string{"2"}, &s.logs),
			InitCommandMock("echo", []string{"3"}, &s.logs),
		},
	}
	s.executionGroup.logger = InitLogger(&LoggerConfig{
		Name: "ExecutionGroupTestSuite",
	})
	s.executionGroup.logger.SetOutput(&s.logs)
}

func (s *ExecutionGroupTestSuite) TestRun() {
	s.executionGroup.Run()
	assert.Contains(s.T(), s.logs.String(), "starting execution group")
	assert.Contains(s.T(), s.logs.String(), "[echo] [1]")
	assert.Contains(s.T(), s.logs.String(), "[echo] [2]")
	assert.Contains(s.T(), s.logs.String(), "[echo] [3]")
	assert.Contains(s.T(), s.logs.String(), "terminated execution group")
}

func (s *ExecutionGroupTestSuite) Test_addPid() {
	expectedPid := 24601
	s.executionGroup.pids = make([]int, 0)
	s.executionGroup.addPid(expectedPid)
	assert.Lenf(s.T(), s.executionGroup.pids, 1, "expected pid %v to have been added but it was not", expectedPid)
	s.executionGroup.addPid(expectedPid)
	assert.Lenf(s.T(), s.executionGroup.pids, 1, "the pid %v seems to have been duplicated when it shouldn't", expectedPid)
}

func (s *ExecutionGroupTestSuite) Test_assertCommandIsValid() {
	// we are running using `go` so there's no reason why it should be unavailable
	expectedApplication := "go"
	testCommand := InitCommand(&CommandConfig{
		Application: expectedApplication,
		Arguments:   []string{},
	})
	s.executionGroup.assertCommandIsValid(testCommand)
}

func (s *ExecutionGroupTestSuite) Test_getExitMessage() {
	expectedPid := 65535
	testCommand := InitCommandMock("test", []string{}, &s.logs)
	exitMessage := s.executionGroup.getExitMessage(testCommand, expectedPid)
	assert.Contains(s.T(), exitMessage, "pid:65535")
	assert.Contains(s.T(), exitMessage, "exit status 0")
}

func (s *ExecutionGroupTestSuite) Test_getStartMessage() {
	expectedPid := 65535
	testCommand := InitCommandMock("test", []string{}, &s.logs)
	exitMessage := s.executionGroup.getStartMessage(testCommand, expectedPid)
	assert.Contains(s.T(), exitMessage, "pid:65535")
	assert.Contains(s.T(), exitMessage, "[command mock]")
}

func (s *ExecutionGroupTestSuite) Test_provisionCommand() {
	expectedPid := 65535
	testCommand := InitCommandMock("test", []string{}, &s.logs)
	s.executionGroup.provisionCommand(testCommand)
	assert.NotNil(s.T(), testCommand.onStart)
	assert.NotNil(s.T(), testCommand.onExit)

	startMessage := testCommand.onStart(expectedPid)
	assert.Contains(s.T(), startMessage, "pid:65535")
	assert.Contains(s.T(), s.executionGroup.pids, expectedPid)

	s.executionGroup.waitGroup.Add(1)
	exitMessage := testCommand.onExit(expectedPid)
	assert.Contains(s.T(), exitMessage, "pid:65535")
	assert.NotContains(s.T(), s.executionGroup.pids, expectedPid)

}

func (s *ExecutionGroupTestSuite) Test_removePid() {
	expectedPid := 24601
	s.executionGroup.pids = make([]int, 0)
	s.executionGroup.addPid(expectedPid)
	assert.Lenf(s.T(), s.executionGroup.pids, 1, "expected pid %v to have been added but it was not", expectedPid)
	s.executionGroup.removePid(expectedPid)
	assert.Lenf(s.T(), s.executionGroup.pids, 0, "the pid %v seems to not have been deleted", expectedPid)
}
