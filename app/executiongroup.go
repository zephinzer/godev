package main

import (
	"sync"
)

var ExecutionGroupCount = 0

// ExecutionGroup runs all commands in parallel
type ExecutionGroup struct {
	commands  []ICommand
	waitGroup sync.WaitGroup
	logger    *Logger
}

// Run starts the execution group's commands in parallel
// and waits for all of them to exit
func (executionGroup *ExecutionGroup) Run() {
	ExecutionGroupCount++
	defer executionGroup.logger.Debugf("execution group[%v] exited", ExecutionGroupCount)
	executionGroup.logger.Debugf("execution group[%v] is starting", ExecutionGroupCount)
	for _, command := range executionGroup.commands {
		if err := command.IsValid(); err != nil {
			executionGroup.logger.Error(err)
		} else {
			executionGroup.waitGroup.Add(1)
			go func() {
				for {
					select {
					case err := <-*command.GetStatus(): // Command letting us know its done
						executionGroup.handleCommandStatus(command, err)
					default:
					}
				}
			}()
			executionGroup.logger.Tracef("command[%s] is starting", command.GetID())
			go command.Run()
		}
	}
	executionGroup.logger.Tracef("waiting for commands to complete running")
	executionGroup.waitGroup.Wait()
}

func (executionGroup *ExecutionGroup) handleCommandStatus(command ICommand, err error) {
	if err != nil {
		executionGroup.logger.Warnf("command[%s] exited with: %s", command.GetID(), err)
	} else {
		executionGroup.logger.Debugf("command[%s] exited without error", command.GetID())
	}
	executionGroup.waitGroup.Done()
}

func (executionGroup *ExecutionGroup) Terminate() {
	for _, command := range executionGroup.commands {
		if command.IsRunning() {
			executionGroup.logger.Tracef("sending SIGINT to command %v", command.GetID())
			command.Interrupt()
			executionGroup.logger.Tracef("SIGINT sent to command %v", command.GetID())
		}
	}
}

func (executionGroup *ExecutionGroup) IsRunning() bool {
	for _, command := range executionGroup.commands {
		if command.IsRunning() {
			return true
		}
	}
	return false
}
