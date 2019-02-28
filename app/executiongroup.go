package main

import (
	"fmt"
	"os/exec"
	"sync"
)

// ExecutionGroup runs all commands in parallel
type ExecutionGroup struct {
	commands  []ICommand
	waitGroup sync.WaitGroup
	logger    *Logger
	pids      []int
	started   bool
}

// Run starts the execution group's commands in parallel
// and waits for all of them to exit
func (executionGroup *ExecutionGroup) Run() {
	defer executionGroup.logger.Tracef("terminated execution group")
	executionGroup.logger.Tracef("starting execution group")
	executionGroup.started = true
	for _, command := range executionGroup.commands {
		executionGroup.waitGroup.Add(1)
		executionGroup.assertCommandIsValid(command)
		executionGroup.provisionCommand(command)
		go command.Run()
	}
	if len(executionGroup.commands) > 0 {
		executionGroup.waitGroup.Wait()
	}
}

// addPid records the provided :pidToAdd, indicating that
// process is running
func (executionGroup *ExecutionGroup) addPid(pidToAdd int) {
	for _, pid := range executionGroup.pids {
		if pid == pidToAdd {
			return
		}
	}
	executionGroup.pids = append(executionGroup.pids, pidToAdd)
}

// assertCommandIsValid does some sanity checks on the provided
// application before we try to run it
func (executionGroup *ExecutionGroup) assertCommandIsValid(command ICommand) {
	application := command.getApplication()
	if len(application) == 0 {
		panic("application field of command is not specified")
	} else if _, err := exec.LookPath(application); err != nil {
		panic(fmt.Sprintf("application '%s' could not be found: %v", application, err))
	}
}

// getExitMessage encapsulates the process exit message more nicely
func (executionGroup *ExecutionGroup) getExitMessage(command ICommand, pid int) string {
	return fmt.Sprintf(
		"%s\n■ ---------------------------------- pid:%v ■",
		command.getCommand().ProcessState.String(),
		pid,
	)
}

// getExitMessage encapsulates the process start message more nicely
func (executionGroup *ExecutionGroup) getStartMessage(command ICommand, pid int) string {
	return fmt.Sprintf(
		"%v\n► ---------------------------------- pid:%v ►",
		command.getArguments(),
		pid,
	)
}

// provisionCommand initialises the provided :command
func (executionGroup *ExecutionGroup) provisionCommand(command ICommand) {
	command.
		setOnStart(func(pid int) string {
			executionGroup.addPid(pid)
			return executionGroup.getStartMessage(command, pid)
		}).
		setOnExit(func(pid int) string {
			defer executionGroup.waitGroup.Done()
			executionGroup.removePid(pid)
			return executionGroup.getExitMessage(command, pid)
		})
}

// removePid unrecords the provided :pidToRemove
func (executionGroup *ExecutionGroup) removePid(pidToRemove int) {
	for index, pid := range executionGroup.pids {
		if pid == pidToRemove {
			executionGroup.pids = append(
				executionGroup.pids[:index],
				executionGroup.pids[index+1:]...,
			)
			break
		}
	}
}
