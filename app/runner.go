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
		logger: InitLogger(&LoggerConfig{Name: "runner", Format: "production", Level: config.LogLevel}),
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
}

func (runner *Runner) Trigger() {
	RunnerTriggerCount++
	defer runner.logger.Tracef("terminated run %v", RunnerTriggerCount)
	runner.logger.Tracef("initialising run %v", RunnerTriggerCount)
	runner.terminateExistingPipeline()
	runner.waitGroup.Add(1)
	go runner.startPipeline()
	runner.waitGroup.Wait()
}

func (runner *Runner) terminateExistingPipeline() {
	defer func() {
		if r := recover(); r != nil {
			runner.logger.Warn(r)
		}
	}()
	runner.waitGroup.Done()
}
