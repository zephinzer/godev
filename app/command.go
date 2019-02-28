package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

type ICommand interface {
	Run()
	getApplication() string
	getArguments() []string
	getCommand() *exec.Cmd
	setOnStart(func(int) string) ICommand
	setOnExit(func(int) string) ICommand
}

func InitCommand(config *CommandConfig) *Command {
	command := &Command{}
	command.config = config
	command.logger = InitLogger(&LoggerConfig{
		Name:   "command",
		Format: "production",
		Level:  config.LogLevel,
		AdditionalFields: &map[string]interface{}{
			"submodule": path.Base(fmt.Sprintf("%s", config.Application)),
		},
	})
	return command
}

type CommandConfig struct {
	Application string
	Arguments   []string
	LogLevel    LogLevel
}

// Command is the atomic command to run
type Command struct {
	config  *CommandConfig
	cmd     *exec.Cmd
	logger  *Logger
	onStart func(int) string
	onExit  func(int) string
}

// Run executes the command
func (command *Command) Run() {
	if command.onStart == nil {
		command.onStart = func(int) string { return "" }
	}
	if command.onExit == nil {
		command.onExit = func(int) string { return "" }
	}
	command.cmd = exec.Command(
		command.config.Application,
		command.config.Arguments...,
	)
	command.cmd.Stderr = os.Stderr
	command.cmd.Stdout = os.Stdout
	pidReported := false
	go func() {
		for {
			if !pidReported {
				if command.cmd.Process != nil && command.cmd.Process.Pid > 0 {
					command.logger.Info(command.onStart(command.cmd.Process.Pid))
					pidReported = true
				}
			} else if command.cmd.ProcessState != nil && command.cmd.ProcessState.Exited() {
				command.logger.Info(command.onExit(command.cmd.ProcessState.Pid()))
				break
			}
		}
	}()

	runErr := command.cmd.Run()
	if runErr != nil {
		command.logger.Warn(runErr)
	}
}

func (command *Command) getApplication() string {
	return command.config.Application
}

func (command *Command) getArguments() []string {
	return command.config.Arguments
}

func (command *Command) getCommand() *exec.Cmd {
	return command.cmd
}

func (command *Command) setOnStart(handler func(int) string) ICommand {
	command.onStart = handler
	return command
}

func (command *Command) setOnExit(handler func(int) string) ICommand {
	command.onExit = handler
	return command
}

func (command *Command) setLogger(logger *Logger) *Command {
	command.logger = logger
	return command
}
