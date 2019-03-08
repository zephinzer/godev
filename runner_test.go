package main

import (
	"bytes"
	"sync"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RunnerTestSuite struct {
	suite.Suite
	runner               *Runner
	logs                 bytes.Buffer
	executionGroupLogger *Logger
}

func TestRunnerTestSuite(t *testing.T) {
	suite.Run(t, new(RunnerTestSuite))
}

func (s *RunnerTestSuite) SetupTest() {
	logger := InitLogger(&LoggerConfig{
		Name: "TestRunnerTestSuite",
	})
	logger.SetOutput(&s.logs)
	executionGroups := []*ExecutionGroup{
		&ExecutionGroup{
			commands: []*Command{
				mockCommand("echo", []string{"runner 1.0"}, &s.logs),
				mockCommand("echo", []string{"runner 1.1"}, &s.logs),
			},
			logger: logger,
		},
		&ExecutionGroup{
			commands: []*Command{
				mockCommand("echo", []string{"runner 2"}, &s.logs),
			},
			logger: logger,
		},
	}
	s.runner = InitRunner(&RunnerConfig{
		Pipeline: executionGroups,
		LogLevel: "trace",
	})
	s.runner.logger.SetOutput(&s.logs)
}

func (s *RunnerTestSuite) Test_startPipeline() {
	s.runner.startPipeline()
	assert.Contains(s.T(), s.logs.String(), "starting pipeline")
	assert.Contains(s.T(), s.logs.String(), "completed pipeline")
}

func (s *RunnerTestSuite) Test_terminateIfRunning_withoutRunningCommand() {
	s.runner.terminateIfRunning()
	assert.Contains(s.T(), s.logs.String(), "is not running")
}

func (s *RunnerTestSuite) Test_terminateIfRunning_withRunningCommand() {
	commandOfInterest := s.runner.config.Pipeline[0].commands[0]
	commandOfInterest.started = true
	commandOfInterest.stopped = false
	commandOfInterest = s.runner.config.Pipeline[0].commands[1]
	commandOfInterest.started = false
	commandOfInterest.stopped = false
	commandOfInterest = s.runner.config.Pipeline[1].commands[0]
	commandOfInterest.started = false
	commandOfInterest.stopped = false
	commandOfInterest = s.runner.config.Pipeline[0].commands[0]
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		select {
		case signal := <-commandOfInterest.signal:
			assert.Equal(s.T(), syscall.SIGINT, signal)
			wg.Done()
			return
		}
	}()
	s.runner.terminateIfRunning()
	wg.Wait()
	assert.Contains(s.T(), s.logs.String(), "terminating pipeline")
	assert.Contains(s.T(), s.logs.String(), "sending SIGINT to command")
	assert.Contains(s.T(), s.logs.String(), "SIGINT received by command")
	assert.Contains(s.T(), s.logs.String(), "SIGINT sent to command")
	assert.Contains(s.T(), s.logs.String(), "terminated pipeline")
}
