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

// Runner is the main component responsible for running the execution pipeline
type Runner struct {
	config    *RunnerConfig
	logger    *Logger
	waitGroup sync.WaitGroup
	started   bool
	stopped   bool
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
		started: false,
		stopped: false,
	}
	return runner
}

func (runner *Runner) startPipeline() {
	RunnerTriggerCount++
	defer runner.logger.Tracef("completed pipeline %v", RunnerTriggerCount)
	runner.logger.Tracef("starting pipeline %v", RunnerTriggerCount)
	executionGroupCount := len(runner.config.Pipeline)
	runner.started = true
	for index, executionGroup := range runner.config.Pipeline {
		executionGroup.logger = InitLogger(&LoggerConfig{
			Name:   "run",
			Format: "production",
			Level:  runner.config.LogLevel,
			AdditionalFields: &map[string]interface{}{
				"submodule": fmt.Sprintf("%v/%v/%v]", RunnerTriggerCount, index+1, executionGroupCount),
			},
		})
		executionGroup.Run()
	}
	runner.stopped = true
}

// Trigger triggers the pipeline
func (runner *Runner) Trigger() {
	runner.started = false
	runner.stopped = false
	runner.terminateIfRunning()
	go runner.startPipeline()
}

func (runner *Runner) terminateIfRunning() {
	defer func() {
		if r := recover(); r != nil {
			runner.logger.Warn(r)
		}
	}()
	for index, executionGroup := range runner.config.Pipeline {
		if executionGroup.IsRunning() {
			runner.logger.Infof("terminating pipeline %v...", RunnerTriggerCount)
			executionGroup.Terminate()
			runner.logger.Infof("terminated pipeline %v", RunnerTriggerCount)
		} else {
			runner.logger.Tracef("execution group %v/%v is not running", index, len(runner.config.Pipeline))
		}
	}
}
