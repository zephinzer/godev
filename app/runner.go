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

var RunnerTriggerCount = 0

type Runner struct {
	config    *RunnerConfig
	logger    *Logger
	waitGroup sync.WaitGroup
}

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
	executionGroupCount := len(runner.config.Pipeline)
	for index, executionGroup := range runner.config.Pipeline {
		executionGroup.logger = InitLogger(&LoggerConfig{
			Name:   fmt.Sprintf("run[%v]", RunnerTriggerCount),
			Format: "production",
			Level:  runner.config.LogLevel,
			AdditionalFields: &map[string]interface{}{
				"submodule": fmt.Sprintf("group[%v/%v]", index+1, executionGroupCount),
			},
		})
		executionGroup.Run()
	}
	runner.waitGroup.Done()
}

func (runner *Runner) isPipelineRunning() bool {
	for _, executionGroup := range runner.config.Pipeline {
		if executionGroup.started && len(executionGroup.pids) > 0 {
			return true
		}
	}
	return false
}

func (runner *Runner) Trigger() {
	RunnerTriggerCount++
	defer runner.logger.Tracef("completed pipeline %v", RunnerTriggerCount)
	runner.logger.Tracef("starting pipeline %v", RunnerTriggerCount)
	runner.waitGroup.Add(1)
	runner.terminateExistingPipeline()
	go runner.startPipeline()
	runner.waitGroup.Wait()
}

func (runner *Runner) terminateExistingPipeline() {
	defer func() {
		if r := recover(); r != nil {
			runner.logger.Warn(r)
		}
	}()
	if runner.isPipelineRunning() {
		runner.waitGroup.Done()
	}
}
