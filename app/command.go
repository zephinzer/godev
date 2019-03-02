package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"
)

var CommandCount = 0

// ICommand is the interface for the Command class
type ICommand interface {
	Run()
	getCommand() *exec.Cmd
	getID() string
	getStatus() *chan error
	IsRunning() bool
	IsValid() error
	Interrupt()
	Terminate()
	Kill()
}

// InitCommand is for creating a new Command
func InitCommand(config *CommandConfig) *Command {
	CommandCount++
	command := &Command{
		id: CommandCount,
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
	LogLevel    LogLevel
}

// Command is the atomic command to run
type Command struct {
	id         int
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

// Interrupt sends SIGINT to the command
func (command *Command) Interrupt() {
	command.logger.Tracef("SIGINT received by command %s", command.getID())
	command.logger.Tracef("command[%v] status: %v/%v, msg: SIGINT >>> %v", command.getID(), command.started, command.terminated, &command.signal)
	command.signal <- syscall.SIGINT
}

// Terminate sends SIGTERM to the command
func (command *Command) Terminate() {
	command.logger.Tracef("SIGTERM received by command %s", command.getID())
	command.logger.Tracef("command[%v] status: %v/%v, msg: SIGTERM >>> %v", command.getID(), command.started, command.terminated, &command.signal)
	command.signal <- syscall.SIGTERM
}

// Kill sends SIGKILL to the command
func (command *Command) Kill() {
	command.logger.Tracef("SIGKILL received by command %s", command.getID())
	command.logger.Tracef("command[%v] status: %v/%v, msg: SIGKILL >>> %v", command.getID(), command.started, command.terminated, &command.signal)
	command.signal <- syscall.SIGKILL
}

// Run executes the command
func (command *Command) Run() {
	command.logger.Tracef("command[%s] is starting", command.getID())
	command.initialise()
	go command.handleProcessRun()
	go func() error {
		for {
			select {
			case signal := <-command.signal: // caller -> Command: shut down please
				command.handleSignalReceived(signal)
				return nil
			case cmdRunStatus := <-command.run: // process -> Command: i'm done here
				return command.handleProcessExited(cmdRunStatus)
			default: // just run
				command.handleRunning()
			}
		}
	}()
	select {
	case terminateCommand := <-command.terminated:
		command.handleTerminated(terminateCommand)
	}
}

func (command *Command) handleProcessExited(status error) error {
	command.logger.Tracef("process status: %v", command.cmd.ProcessState)
	command.terminated <- status
	return nil
}

func (command *Command) handleProcessRun() {
	command.started = true
	command.run <- command.cmd.Run()
}

func (command *Command) handleRunning() {
	if !command.reported {
		if command.cmd.Process != nil {
			command.logger.Infof("%v\n► ─────────────────────────────── pid:%v ►", strings.Join(command.config.Arguments, ","), command.cmd.Process.Pid)
			command.reported = true
		}
	}
}

func (command *Command) handleSignalReceived(signal os.Signal) error {
	command.logger.Tracef("caller sent signal %v", signal)
	command.terminated <- errors.New(signal.String())
	if err := command.cmd.Process.Signal(signal); err != nil {
		command.logger.Warn(err)
		return err
	}
	return nil
}

func (command *Command) handleTerminated(terminateCommand error) {
	command.logger.Tracef("command[%s] is exiting (%v)", command.getID(), terminateCommand)
	command.logger.Infof("%v\n■ ─────────────────────────────── pid:%v ■", strings.Join(command.config.Arguments, ","), command.cmd.Process.Pid)
	command.stopped = true
	command.status <- terminateCommand
}

func (command *Command) initialise() {
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
	command.cmd.Stderr = os.Stderr
	command.cmd.Stdout = os.Stdout
}

func (command *Command) getCommand() *exec.Cmd {
	return command.cmd
}

func (command *Command) getID() string {
	return strconv.Itoa(command.id)
}

func (command *Command) getStatus() *chan error {
	return &command.status
}

func (command *Command) IsRunning() bool {
	command.logger.Infof("command %v: (s/t) %v/%v", command.id, command.started, command.terminated)
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

func (command *Command) setLogger(logger *Logger) *Command {
	command.logger = logger
	return command
}
