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
	defer executionGroup.logger.Tracef("terminated execution group")
	executionGroup.logger.Tracef("starting execution group")
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
			command.logger.Infof("%v\n► ---------------------------------- pid:%v ►", command.arguments, pid)
		}
		command.onExit = func(pid int) {
			command.logger.Infof("%s\n■ ---------------------------------- pid:%v ■", command.cmd.ProcessState.String(), pid)
			executionGroup.waitGroup.Done()
		}
		go command.Run()
	}
	executionGroup.waitGroup.Wait()
}
