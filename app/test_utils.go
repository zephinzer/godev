package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createFile(t *testing.T, pathToFile string) {
	testFileCreation, err := os.Create(pathToFile)
	assert.Nilf(
		t,
		err,
		"unexpected error when creating file at '%s': %s",
		pathToFile,
		err,
	)
	defer testFileCreation.Close()
}

type MockCommand struct {
	Command
}

func mockCommand(application string, arguments []string, logOutput *bytes.Buffer) *Command {
	command := &Command{
		id: fmt.Sprintf("%s%v", application, arguments),
		config: &CommandConfig{
			Application: application,
			Arguments:   arguments,
		},
		logger: InitLogger(&LoggerConfig{
			Name:   application,
			Format: "production",
			Level:  "trace",
		}),
	}
	command.logger.SetOutput(logOutput)
	return command
}

func expectError(t *testing.T) func() {
	return func() {
		assert.NotNil(t, recover(), "expected an error but none was panicked")
	}
}

func expectNoError(t *testing.T) func() {
	return func() {
		err := recover()
		assert.Nilf(t, err, "expected no errors but '%s' was panicked", err)
	}
}

func removeFile(t *testing.T, pathToFile string) {
	err := os.Remove(pathToFile)
	assert.Nilf(
		t,
		err,
		"unexpected error when removing file at '%s': %s",
		pathToFile,
		err,
	)
}
