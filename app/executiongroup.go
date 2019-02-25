package main

import (
	"fmt"
	"strings"
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
		executionGroup.waitGroup.Add(1)
		command.logger = InitLogger(&LoggerConfig{
			Name:   "command",
			Format: "production",
			Level:  "trace",
			AdditionalFields: &map[string]interface{}{
				"app":  fmt.Sprintf("%s", command.application),
				"args": fmt.Sprintf("%s", strings.Join(command.arguments, " ")),
			},
		})
		command.onStart = func(pid int) {
			command.logger.Infof("[START] ----- pid:%v -----", pid)
		}
		command.onExit = func(pid int) {
			command.logger.Infof("[STOP]  ----- pid:%v -----", pid)
			executionGroup.waitGroup.Done()
		}
		go command.Run()
	}
	executionGroup.waitGroup.Wait()
}
