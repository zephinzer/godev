package main

import (
	"bytes"
	"errors"
	"os/exec"
	"path"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CommandTestSuite struct {
	suite.Suite
}

func TestCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CommandTestSuite))
}

func (s *CommandTestSuite) TestGetID() {
	expectedID := "id"
	cmd := &Command{id: expectedID}
	assert.Equal(s.T(), cmd.GetID(), expectedID)
}

func (s *CommandTestSuite) TestGetStatus() {
	cmd := &Command{status: make(chan error, 0)}
	statusRef := cmd.GetStatus()
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wait <-chan time.Time) {
		select {
		case <-wait:
			cmd.status <- errors.New("testGetStatus")
		}
	}(time.After(100 * time.Millisecond))
	go func(status *chan error) {
		select {
		case status := <-*status:
			assert.Equal(s.T(), "testGetStatus", status.Error())
			wg.Done()
			return
		}
	}(statusRef)
	wg.Wait()
}

func (s *CommandTestSuite) TestIsRunning() {
	cmd := &Command{started: false, stopped: false}
	assert.False(s.T(), cmd.IsRunning())
	cmd.started = true
	assert.True(s.T(), cmd.IsRunning())
	cmd.stopped = true
	assert.False(s.T(), cmd.IsRunning())
}

func (s *CommandTestSuite) TestIsValidFromRegisteredPath() {
	// we are running using `go` so there's no reason why it should be unavailable
	expectedApplication := "go"
	testCommand := InitCommand(&CommandConfig{
		Application: expectedApplication,
		Arguments:   []string{},
	})
	err := testCommand.IsValid()
	assert.Nil(s.T(), err)
}

func (s *CommandTestSuite) TestIsValidFromAbsolutePathNoPermissions() {
	cwd := getCurrentWorkingDirectory()
	expectedApplication := path.Join(cwd, "/data/test-exec/nonexec.sh")
	testCommand := InitCommand(&CommandConfig{
		Application: expectedApplication,
		Arguments:   []string{},
	})
	err := testCommand.IsValid()
	assert.NotNil(s.T(), err)
	assert.Contains(s.T(), err.Error(), "you don't have permissions to execute")
}

func (s *CommandTestSuite) TestIsValidFromAbsolutePathWithPermissions() {
	cwd := getCurrentWorkingDirectory()
	expectedApplication := path.Join(cwd, "/data/test-exec/exec.sh")
	testCommand := InitCommand(&CommandConfig{
		Application: expectedApplication,
		Arguments:   []string{},
	})
	err := testCommand.IsValid()
	assert.Nil(s.T(), err)
}

func (s *CommandTestSuite) TestRun() {
	panic("TODO")
}

func (s *CommandTestSuite) TestSendInterrupt() {
	panic("TODO")
}

func (s *CommandTestSuite) Test_handleInitialisation() {
	panic("TODO")
}

func (s *CommandTestSuite) Test_handleProcessExited() {
	panic("TODO")
}

func (s *CommandTestSuite) Test_handleProcessLifecycle() {
	panic("TODO")
}

func (s *CommandTestSuite) Test_handleProcessReporting() {
	panic("TODO")
}

func (s *CommandTestSuite) Test_handleSignalReceived() {
	panic("TODO")
}

func (s *CommandTestSuite) Test_handleStart() {
	panic("TODO")
}

func (s *CommandTestSuite) Test_handleStopped() {
	var logs bytes.Buffer
	var wg sync.WaitGroup
	logger := InitLogger(&LoggerConfig{
		Name:   "cmd",
		Format: "production",
		Level:  "trace",
	})
	logger.SetOutput(&logs)
	cmd := &Command{
		id:      "id",
		logger:  logger,
		stopped: false,
		status:  make(chan error, 0),
		cmd:     exec.Command("test", []string{"b"}...),
	}
	wg.Add(1)
	go func() {
		select {
		case status := <-cmd.status:
			assert.True(s.T(), cmd.stopped)
			assert.Equal(s.T(), status.Error(), "Test_handleStopped")
			wg.Done()
			return
		}
	}()
	cmd.handleStopped(errors.New("Test_handleStopped"))
	wg.Wait()
}
