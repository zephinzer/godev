package main

import (
	"sync"
)

// ExecutionGroupCount keeps track of the execution group count for
// display in the verbose logs - helps to differentiate between
// the different execution groups
var ExecutionGroupCount = 0

// ExecutionGroup runs all commands in parallel
type ExecutionGroup struct {
	commands  []*Command
	waitGroup sync.WaitGroup
	logger    *Logger
}

// IsRunning is for the Runner to check if the execution group
// is still running
func (executionGroup *ExecutionGroup) IsRunning() bool {
	for _, command := range executionGroup.commands {
		if command.IsRunning() {
			return true
		}
	}
	return false
}

// Run starts the execution group's commands in parallel
// and waits for all of them to exit
func (executionGroup *ExecutionGroup) Run() {
	ExecutionGroupCount++
	defer executionGroup.logger.Debugf("execution group[%v] exited", ExecutionGroupCount)
	executionGroup.logger.Debugf("execution group[%v] is starting...", ExecutionGroupCount)
	for _, command := range executionGroup.commands {
		if err := command.IsValid(); err != nil {
			executionGroup.logger.Error(err)
		} else {
			go func(commandStatus *chan error) {
				for {
					select {
					case err := <-*commandStatus: // Command letting us know its done
						executionGroup.handleCommandStatus(command, err)
						return
					default:
					}
				}
			}(command.GetStatus())
			executionGroup.logger.Tracef("command[%s] is starting", command.GetID())
			executionGroup.waitGroup.Add(1)
			go command.Run()
		}
	}
	executionGroup.logger.Tracef("waiting for commands to complete running...")
	executionGroup.waitGroup.Wait()
}

// Terminate terminates this instance of the execution group, used when
// the Runner receives a signal to start a new pipeline
func (executionGroup *ExecutionGroup) Terminate() {
	for _, command := range executionGroup.commands {
		if command.IsRunning() {
			executionGroup.logger.Tracef("sending SIGINT to command %v", command.GetID())
			command.SendInterrupt()
			executionGroup.logger.Tracef("SIGINT sent to command %v", command.GetID())
		}
	}
}

func (executionGroup *ExecutionGroup) handleCommandStatus(command *Command, err error) {
	defer func() {
		if r := recover(); r != nil {
			executionGroup.logger.Warn(r)
		}
	}()
	if err != nil {
		executionGroup.logger.Warnf("command[%s] exited with: %s", command.GetID(), err)
	} else {
		executionGroup.logger.Debugf("command[%s] exited without error", command.GetID())
	}
	executionGroup.waitGroup.Done()
}
