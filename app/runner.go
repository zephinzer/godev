package main

import (
	"fmt"
	"sync"
)

// RunnerConfig configures the Runner
type RunnerConfig struct {
	pipeline  []*ExecutionGroup
	waitGroup sync.WaitGroup
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
		logger: InitLogger(&LoggerConfig{Name: "runner", Format: "production", Level: "trace"}),
	}
	return runner
}

func (runner *Runner) startPipeline() {
	executionGroupCount := len(runner.config.pipeline)
	for index, executionGroup := range runner.config.pipeline {
		executionGroup.logger = InitLogger(&LoggerConfig{
			Name:   "execution_group",
			Format: "production",
			Level:  "trace",
			AdditionalFields: &map[string]interface{}{
				"id":    RunnerTriggerCount,
				"group": fmt.Sprintf("%v/%v", index+1, executionGroupCount),
			},
		})
		executionGroup.Run()
	}
}

func (runner *Runner) Trigger() {
	RunnerTriggerCount++
	defer runner.logger.Infof("terminated run %v", RunnerTriggerCount)
	runner.logger.Infof("initialising run %v", RunnerTriggerCount)
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
