package main

import (
	"os"
	"os/exec"
)

// Command is the atomic command to run
type Command struct {
	application string
	arguments   []string
	cmd         *exec.Cmd
	logger      *Logger
	onStart     func(int)
	onExit      func(int)
}

// Run executes the command
func (command *Command) Run() {
	if command.onStart == nil {
		command.onStart = func(int) {}
	}
	if command.onExit == nil {
		command.onExit = func(int) {}
	}
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
		command.logger.Warn(runErr)
	}
}
