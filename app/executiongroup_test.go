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
