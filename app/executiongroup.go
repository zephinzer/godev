package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
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
		if valid, err := executionGroup.isCommandValid(command); !valid {
			executionGroup.logger.Error(err)
		} else {
			executionGroup.provisionCommand(command)
			go command.Run()
		}
	}
	if len(executionGroup.commands) > 0 {
		executionGroup.logger.Tracef("now waiting for commands to complete running")
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

// isCommandValid does some sanity checks on the provided
// application before we try to run it
func (executionGroup *ExecutionGroup) isCommandValid(command ICommand) (bool, error) {
	application := command.getApplication()
	if len(application) == 0 {
		return false, errors.New("application field of command is not specified")
	} else if path.IsAbs(application) {
		_, err := os.Lstat(application)
		if err != nil && os.IsNotExist(err) {
			return false, fmt.Errorf("application at '%s' could not be found", application)
		} else if err != nil {
			return false, err
		}
	}
	if _, err := exec.LookPath(application); err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			return false, fmt.Errorf("looks like you don't have permissions to execute '%s', run 'chmod +x %s' and try again", application, application)
		}
		return false, fmt.Errorf("error occurred while running application '%s': %v", application, err)
	}
	return true, nil
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
	executionGroup.waitGroup.Add(1)
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
