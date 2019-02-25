package main

import (
	"fmt"
	"sync"
)

// ExecutionGroup runs all commands in parallel
type ExecutionGroup struct {
	commands  []*Command
	waitGroup sync.WaitGroup
	logger    *Logger
	pids      []int
}

func (executionGroup *ExecutionGroup) Run() {
	defer executionGroup.logger.Info("terminated execution group")
	executionGroup.logger.Info("starting execution group")
	for _, command := range executionGroup.commands {
		command.logger = InitLogger(&LoggerConfig{
			Name:   "command",
			Format: "production",
			Level:  "trace",
			AdditionalFields: &map[string]interface{}{
				"application": fmt.Sprintf("%s", command.application),
				"arguments":   fmt.Sprintf("%s", command.arguments),
			},
		})
		command.onStart = func(pid int) {
			executionGroup.waitGroup.Add(1)
			executionGroup.logger.Infof("process %v has started", pid)
		}
		command.onExit = func(pid int) {
			executionGroup.logger.Infof("process %v has exitted", pid)
			executionGroup.waitGroup.Done()
		}
		go command.Run()
	}
	executionGroup.waitGroup.Wait()
}

// RunnerConfig configures the Runner
type RunnerConfig struct {
	pipeline  []ExecutionGroup
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
	for index, executionGroup := range runner.config.pipeline {
		executionGroup.logger = InitLogger(&LoggerConfig{
			Name:   "execution_group",
			Format: "production",
			Level:  "trace",
			AdditionalFields: &map[string]interface{}{
				"id":    RunnerTriggerCount,
				"index": index,
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
