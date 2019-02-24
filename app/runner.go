package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

// Command is the atomic command to run
type Command struct {
	application string
	arguments   []string
	cmd         *exec.Cmd
	onStart     func(int)
	onExit      func(int)
}

func (command *Command) Run() {
	if command.onStart == nil {
		command.onStart = func(int) {}
	}
	if command.onExit == nil {
		command.onExit = func(int) {}
	}
	log.Println("running command: ", command.application)
	command.cmd = exec.Command(command.application, command.arguments...)

	command.cmd.Stderr = os.Stderr
	command.cmd.Stdout = os.Stdout

	pidReported := false
	go func() {
		for {
			if !pidReported {
				if command.cmd.Process != nil && command.cmd.Process.Pid > 0 {
					command.onStart(command.cmd.Process.Pid)
					pidReported = true
				}
			} else if command.cmd.ProcessState != nil && command.cmd.ProcessState.Exited() {
				command.onExit(command.cmd.ProcessState.Pid())
				break
			}
		}
	}()

	runErr := command.cmd.Run()
	if runErr != nil {
		panic(runErr)
	}
}

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
			Name:   fmt.Sprintf("run[%v/%v]", RunnerTriggerCount, index),
			Format: "production",
			Level:  "trace",
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
