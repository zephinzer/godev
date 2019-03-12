package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CommandTestSuite struct {
	suite.Suite
	command    *Command
	logs       bytes.Buffer
	expectedID string
}

func TestCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CommandTestSuite))
}

func (s *CommandTestSuite) SetupTest() {
	s.expectedID = "CommandTestSuiteCommandID"
	logger := InitLogger(&LoggerConfig{
		Name:   "CommandTestSuite",
		Format: "production",
		Level:  "trace",
	})
	logger.SetOutput(&s.logs)
	config := &CommandConfig{
		Application: "go",
		Arguments:   []string{"version"},
	}
	s.command = &Command{
		id:     s.expectedID,
		config: config,
		logger: logger,
	}
	s.command.handleInitialisation()
}

func (s *CommandTestSuite) TestGetID() {
	assert.Equal(s.T(), s.command.GetID(), s.expectedID)
}

func (s *CommandTestSuite) TestGetStatus() {
	statusRef := s.command.GetStatus()
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wait <-chan time.Time) {
		select {
		case <-wait:
			s.command.status <- errors.New("testGetStatus")
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
	t := s.T()
	s.command.started = false
	s.command.stopped = false
	assert.False(t, s.command.IsRunning())
	s.command.started = true
	assert.True(t, s.command.IsRunning())
	s.command.stopped = true
	assert.False(t, s.command.IsRunning())
}

func (s *CommandTestSuite) TestIsValid_FromRegisteredPath() {
	err := s.command.IsValid()
	assert.Nil(s.T(), err)
}

func (s *CommandTestSuite) TestIsValid_FromAbsolutePath_NoPermissions() {
	cwd := getCurrentWorkingDirectory()
	s.command.config.Application = path.Join(cwd, "/data/test-exec/nonexec.sh")
	s.command.handleInitialisation()
	err := s.command.IsValid()
	assert.NotNil(s.T(), err)
	assert.Contains(s.T(), err.Error(), "permission denied")
}

func (s *CommandTestSuite) TestIsValid_FromAbsolutePath_WithPermissions() {
	cwd := getCurrentWorkingDirectory()
	s.command.config.Application = path.Join(cwd, "/data/test-exec/exec.sh")
	s.command.handleInitialisation()
	err := s.command.IsValid()
	assert.Nil(s.T(), err)
}

func (s *CommandTestSuite) TestRun() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			select {
			case status := <-s.command.status:
				assert.Equal(s.T(), nil, status)
				assert.Contains(s.T(), s.logs.String(), "command[CommandTestSuiteCommandID] is starting")
				wg.Done()
				return
			default:
			}
		}
	}()
	s.command.Run()
	wg.Wait()
}

func (s *CommandTestSuite) TestSendInterrupt() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		select {
		case signal := <-s.command.signal:
			assert.Equal(s.T(), "interrupt", signal.String())
			assert.Contains(s.T(), s.logs.String(), "SIGINT received")
			wg.Done()
			return
		}
	}()
	s.command.SendInterrupt()
	wg.Wait()
}

func (s *CommandTestSuite) Test_handleInitialisation() {
	t := s.T()
	expectedDir := "/some/directory"
	expectedEnv := []string{
		"A=1",
		"B=2",
	}
	cmd := &Command{
		config: &CommandConfig{
			Application: "go",
			Arguments:   []string{"version"},
			Directory:   expectedDir,
			Environment: expectedEnv,
		},
	}
	cmd.handleInitialisation()
	assert.NotNil(t, cmd.signal)
	assert.NotNil(t, cmd.status)
	assert.NotNil(t, cmd.run)
	assert.NotNil(t, cmd.terminated)
	assert.NotNil(t, cmd.cmd.Stderr)
	assert.NotNil(t, cmd.cmd.Stdout)
	assert.Equal(t, expectedDir, cmd.cmd.Dir)
	assert.False(t, cmd.reported)
	assert.False(t, cmd.started)
	assert.False(t, cmd.stopped)
	assert.Len(t, cmd.cmd.Env, len(expectedEnv))
	for index, env := range expectedEnv {
		assert.Equal(t, env, cmd.cmd.Env[index])
	}
}

func (s *CommandTestSuite) Test_handleProcessExited() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		select {
		case status := <-s.command.terminated:
			assert.Equal(s.T(), "Test_handleProcessExited", status.Error())
			wg.Done()
		}
	}()
	s.command.handleProcessExited(errors.New("Test_handleProcessExited"))
	wg.Wait()
}

func (s *CommandTestSuite) Test_handleProcessLifecycleCallerSaysStop() {
	var wg sync.WaitGroup
	s.command.cmd.Process = &os.Process{}
	wg.Add(1)
	go func() {
		select {
		case signal := <-s.command.terminated:
			assert.Equal(s.T(), "interrupt", signal.Error())
			wg.Done()
		}
	}()
	go func(time <-chan time.Time) {
		select {
		case <-time:
			s.command.signal <- syscall.SIGINT
			return
		}
	}(time.After(100 * time.Millisecond))
	go s.command.handleProcessLifecycle()
	wg.Wait()
}

func (s *CommandTestSuite) Test_handleProcessLifecycleProcessSaysStopped() {
	var wg sync.WaitGroup
	s.command.cmd.Process = &os.Process{}
	wg.Add(1)
	go func() {
		select {
		case status := <-s.command.terminated:
			assert.Equal(s.T(), "Test_handleProcessLifecycleProcessSaysStopped", status.Error())
			wg.Done()
		}
	}()
	go func(time <-chan time.Time) {
		select {
		case <-time:
			s.command.run <- errors.New("Test_handleProcessLifecycleProcessSaysStopped")
			return
		}
	}(time.After(100 * time.Millisecond))
	go s.command.handleProcessLifecycle()
	wg.Wait()
}

func (s *CommandTestSuite) Test_handleProcessReporting() {
	s.command.reported = false
	s.command.cmd.Process = &os.Process{Pid: -1}
	s.command.handleProcessReporting()
	assert.True(s.T(), s.command.reported)
	assert.Contains(s.T(), s.logs.String(), "pid:-1 id:CommandTestSuiteCommandID")
}

func (s *CommandTestSuite) Test_handleSignalReceived() {
	sigcalls := []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL}
	sigstring := []string{"interrupt", "terminated", "killed"}
	var wg sync.WaitGroup
	s.command.cmd.Process = &os.Process{}
	for i := 0; i < len(sigcalls); i++ {
		wg.Add(1)
		go func(j int) {
			select {
			case signal := <-s.command.terminated:
				assert.Equal(s.T(), sigstring[j], signal.Error())
				wg.Done()
			}
		}(i)
		s.command.handleSignalReceived(sigcalls[i])
		wg.Wait()
	}
}

func (s *CommandTestSuite) Test_handleStart() {
	t := s.T()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		select {
		case err := <-s.command.run:
			assert.True(t, s.command.started)
			assert.Nil(t, err)
			wg.Done()
			return
		}
	}()
	go s.command.handleStart()
	wg.Wait()
}

func (s *CommandTestSuite) Test_handleStopped() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		select {
		case status := <-s.command.status:
			assert.True(s.T(), s.command.stopped)
			assert.Equal(s.T(), status.Error(), "Test_handleStopped")
			assert.Contains(s.T(), s.logs.String(), fmt.Sprintf("command[%s] is exiting", s.command.id))
			wg.Done()
			return
		}
	}()
	s.command.handleStopped(errors.New("Test_handleStopped"))
	wg.Wait()
}
