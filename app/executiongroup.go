package main

import (
	"fmt"
	"path"
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
				"submodule": path.Base(fmt.Sprintf("%s", command.application)),
			},
		})
		command.onStart = func(pid int) {
			command.logger.Infof("\n► ---------------------------------- pid:%v ▼", pid)
		}
		command.onExit = func(pid int) {
			command.logger.Infof("\n■ ---------------------------------- pid:%v ▲", pid)
			executionGroup.waitGroup.Done()
		}
		go command.Run()
	}
	executionGroup.waitGroup.Wait()
}
