package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli"
)

type CLIViewHandlerTestSuite struct {
	suite.Suite
	mockApp     *cli.App
	mockContext *cli.Context
}

func TestCLIViewHandler(t *testing.T) {
	suite.Run(t, new(CLIViewHandlerTestSuite))
}

func (s *CLIViewHandlerTestSuite) SetupTest() {
	s.mockApp = cli.NewApp()
	s.mockContext = cli.NewContext(s.mockApp, nil, nil)
}

func (s *CLIViewHandlerTestSuite) Test_getViewCommand() {
	config := Config{}
	logger := InitLogger(&LoggerConfig{Name: "getViewCommand", Format: "raw", Level: "trace"})
	command := getViewCommand(&config, logger)
	ensureCLICommand(s.T(), command, []string{"view", "V"}, []cli.Flag(nil))
}

func (s *CLIViewHandlerTestSuite) Test_getViewAction() {
	t := s.T()
	config := Config{}
	logger := InitLogger(&LoggerConfig{Name: "getViewAction", Format: "raw", Level: "trace"})
	s.mockApp.Action = getViewAction(&config, logger)
	_ = s.mockApp.Run([]string{"test-run-version", "dockerfile"})
	assert.True(t, config.RunView)
}
