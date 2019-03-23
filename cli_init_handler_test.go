package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli"
)

type CLIInitHandlerTestSuite struct {
	suite.Suite
	mockApp     *cli.App
	mockContext *cli.Context
}

func TestCLIInitHandler(t *testing.T) {
	suite.Run(t, new(CLIInitHandlerTestSuite))
}

func (s *CLIInitHandlerTestSuite) SetupTest() {
	s.mockApp = cli.NewApp()
	s.mockApp.Flags = getInitFlags()
	s.mockContext = cli.NewContext(s.mockApp, nil, nil)
}

func (s *CLIInitHandlerTestSuite) Test_getInitCommand() {
	config := Config{}
	command := getInitCommand(&config)
	ensureCLICommand(s.T(), command, []string{"init", "i"}, getInitFlags())
}

func (s *CLIDefaultHandlerTestSuite) Test_getInitFlags() {
	ensureCLIFlags(s.T(),
		[]string{
			"dir",
		},
		getInitFlags(),
	)
}

func (s *CLIInitHandlerTestSuite) Test_getInitAction() {
	t := s.T()
	config := Config{}
	s.mockApp.Action = getInitAction(&config)
	if err := s.mockApp.Run([]string{"test-run-init"}); err == nil {
		assert.True(t, config.RunInit)
		assert.Equal(t, getCurrentWorkingDirectory(), config.WorkDirectory)
	} else {
		panic(err)
	}
}
