package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
)

// CommandDelimiter is used when demarcating boundaries between
// functions
const CommandDelimiter = "───────────────────────────────────────────"

// CommandProcessStartSymbol is the fancy symbol we use to denote
// the start of a command
const CommandProcessStartSymbol = "►"

// CommandProcessStopSymbol is the fancy symbol we use to denote
// the end of a command
const CommandProcessStopSymbol = "■"

// ICommand is the interface for the Command class
type ICommand interface {
	// runs the command
	Run()
	// gets the id of the command
	GetID() string
	// get a pointer to the status channel
	GetStatus() *chan error
	// checks if the command is still running
	IsRunning() bool
	// checks if the command is valid
	IsValid() error
	// tells command to exit nicely
	SendInterrupt()
}

// InitCommand is for creating a new Command
func InitCommand(config *CommandConfig) *Command {
	commandHash := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%v", config.Application, config.Arguments))))
	command := &Command{
		id: commandHash[:6],
	}
	command.config = config
	command.logger = InitLogger(&LoggerConfig{
		Name:   "command",
		Format: "production",
		Level:  config.LogLevel,
		AdditionalFields: &map[string]interface{}{
			"submodule": path.Base(fmt.Sprintf("%s", config.Application)),
		},
	})
	return command
}

// CommandConfig configures Command
type CommandConfig struct {
	Application string
	Arguments   []string
	Directory   string
	LogLevel    LogLevel
}

// Command is the atomic command to run
type Command struct {
	id         string
	signal     chan os.Signal
	status     chan error
	run        chan error
	terminated chan error
	config     *CommandConfig
	cmd        *exec.Cmd
	logger     *Logger
	started    bool
	reported   bool
	stopped    bool
}

// GetID returns the command's ID, used for the execution group
// to report the running command
func (command *Command) GetID() string {
	return command.id
}

// GetStatus returns the command's status channel for the execution
// group to know when the command has terminated
func (command *Command) GetStatus() *chan error {
	return &command.status
}

// IsRunning allows callers to check if the command is running,
// the logic is tied into the Run()
func (command *Command) IsRunning() bool {
	return command.started && !command.stopped
}

// IsValid does some sanity checks on the provided
// application before we try to run it
func (command *Command) IsValid() error {
	application := command.config.Application
	if len(application) == 0 {
		return errors.New("no application was specified")
	} else if path.IsAbs(application) {
		_, err := os.Lstat(application)
		if err != nil && os.IsNotExist(err) {
			return fmt.Errorf("application at '%s' could not be found", application)
		} else if err != nil {
			return err
		}
	}
	if _, err := exec.LookPath(application); err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			return fmt.Errorf("looks like you don't have permissions to execute '%s',\n  * run 'chmod +x %s' and try again", path.Base(application), application)
		}
		return fmt.Errorf("unexpected error occurred while running application '%s': %v", application, err)
	}
	return nil
}

// Run executes the command
func (command *Command) Run() {
	command.logger.Tracef("command[%s] is starting", command.id)
	command.handleInitialisation()
	go command.handleStart()
	go command.handleProcessLifecycle()
	select {
	case terminateCommand := <-command.terminated:
		command.handleStopped(terminateCommand)
	}
}

// SendInterrupt sends SIGINT to the command
func (command *Command) SendInterrupt() {
	command.logger.Tracef("SIGINT received by command %s", command.id)
	command.logger.Tracef("command[%v] status: %v/%v, msg: SIGINT >>> %v", command.id, command.started, command.terminated, &command.signal)
	command.signal <- syscall.SIGINT
}

func (command *Command) handleInitialisation() {
	if command.config == nil {
		panic("command.config needs to be defined before initialisation can be done")
	}
	command.signal = make(chan os.Signal, 0)
	command.status = make(chan error, 0)
	command.run = make(chan error, 0)
	command.terminated = make(chan error, 0)
	command.started = false
	command.reported = false
	command.stopped = false
	command.cmd = exec.Command(
		command.config.Application,
		command.config.Arguments...,
	)
	command.cmd.Dir = command.config.Directory
	command.cmd.Stderr = os.Stderr
	command.cmd.Stdout = os.Stdout
}

// handleProcessExited handles the exit status being sent by the process
func (command *Command) handleProcessExited(status error) error {
	command.logger.Tracef("process status: %v", command.cmd.ProcessState)
	command.terminated <- status
	return nil
}

func (command *Command) handleProcessLifecycle() error {
	for {
		select {
		case signal := <-command.signal: // caller -> Command: shut down please
			return command.handleSignalReceived(signal)
		case cmdRunStatus := <-command.run: // process -> Command: i'm done here
			return command.handleProcessExited(cmdRunStatus)
		default: // just run
			command.handleProcessReporting()
		}
	}
}

// handleProcessReporting handles the CLI reporting after a process
// has its PID reported
func (command *Command) handleProcessReporting() {
	if !command.reported {
		if command.cmd.Process != nil {
			command.logger.Infof(
				"'%v'\n%s %s pid:%v id:%s %s",
				strings.Join(command.config.Arguments, "', '"),
				CommandProcessStartSymbol,
				CommandDelimiter,
				command.cmd.Process.Pid,
				command.id,
				CommandProcessStartSymbol,
			)
			command.reported = true
		}
	}
}

// handleSignalReceived handles the signal received by the caller
func (command *Command) handleSignalReceived(signal os.Signal) error {
	command.logger.Tracef("caller sent signal %v", signal)
	command.terminated <- errors.New(signal.String())
	if err := command.cmd.Process.Signal(signal); err != nil {
		command.logger.Warn(err)
		return err
	}
	return nil
}

// handleStart starts the process
func (command *Command) handleStart() {
	command.started = true
	command.run <- command.cmd.Run()
}

// handleStopped processes the end of a command as reported
// by (*exec.Cmd).Run or (*exec.Cmd).Wait
func (command *Command) handleStopped(terminateCommand error) {
	command.logger.Tracef("command[%s] is exiting (%v)", command.id, terminateCommand)
	pid := -1
	if command.cmd.Process != nil {
		pid = command.cmd.Process.Pid
	}
	command.logger.Infof(
		"\n%s %s pid:%v id:%s %s",
		CommandProcessStopSymbol,
		CommandDelimiter,
		pid,
		command.id,
		CommandProcessStopSymbol,
	)
	command.stopped = true
	command.status <- terminateCommand
}
