package main

import (
	"fmt"
	"sync"
)

// RunnerConfig configures the Runner
type RunnerConfig struct {
	Pipeline []*ExecutionGroup
	LogLevel LogLevel
}

// RunnerTriggerCount keeps track of the number of piplines run
var RunnerTriggerCount = 0

// Runner is the main component responsible for running the commands
type Runner struct {
	config    *RunnerConfig
	logger    *Logger
	waitGroup sync.WaitGroup
}

// InitRunner initialises a runner
func InitRunner(config *RunnerConfig) *Runner {
	runner := &Runner{
		config: config,
		logger: InitLogger(&LoggerConfig{
			Name:   "runner",
			Format: "production",
			Level:  config.LogLevel},
		),
	}
	return runner
}

func (runner *Runner) startPipeline() {
	defer runner.logger.Tracef("completed pipeline %v", RunnerTriggerCount)
	executionGroupCount := len(runner.config.Pipeline)
	for index, executionGroup := range runner.config.Pipeline {
		executionGroup.logger = InitLogger(&LoggerConfig{
			Name:   "iteration",
			Format: "production",
			Level:  runner.config.LogLevel,
			AdditionalFields: &map[string]interface{}{
				"submodule": fmt.Sprintf("%v/%v/%v]", RunnerTriggerCount, index+1, executionGroupCount),
			},
		})
		executionGroup.Run()
	}
}

func (runner *Runner) IsRunning() bool {
	for _, executionGroup := range runner.config.Pipeline {
		if executionGroup.IsRunning() {
			return true
		}
	}
	return false
}

func (runner *Runner) Trigger() {
	runner.terminateExistingPipeline()
	RunnerTriggerCount++
	go runner.startPipeline()
}

func (runner *Runner) terminateExistingPipeline() {
	defer func() {
		if r := recover(); r != nil {
			runner.logger.Warn(r)
		}
	}()
	for index, executionGroup := range runner.config.Pipeline {
		if executionGroup.IsRunning() {
			runner.logger.Warnf("terminating pipeline %v...", RunnerTriggerCount)
			executionGroup.Terminate()
			runner.logger.Errorf("terminated pipeline %v", RunnerTriggerCount)
		} else {
			runner.logger.Infof("execution group %v is not running", index)
		}
	}
}
