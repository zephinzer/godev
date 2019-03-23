package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli"
)

type CLIVersionHandlerTestSuite struct {
	suite.Suite
	mockApp     *cli.App
	mockContext *cli.Context
}

func TestCLIVersionHandler(t *testing.T) {
	suite.Run(t, new(CLIVersionHandlerTestSuite))
}

func (s *CLIVersionHandlerTestSuite) SetupTest() {
	s.mockApp = cli.NewApp()
	s.mockApp.Flags = getVersionFlags()
	s.mockContext = cli.NewContext(s.mockApp, nil, nil)
}

func (s *CLIVersionHandlerTestSuite) Test_getVersionCommand() {
	config := Config{}
	logger := InitLogger(&LoggerConfig{Name: "getVersionCommand", Format: "raw", Level: "trace"})
	command := getVersionCommand(&config, logger)
	ensureCLICommand(s.T(), command, []string{"version", "v"}, getVersionFlags())
}

func (s *CLIVersionHandlerTestSuite) Test_getVersionFlags() {
	ensureCLIFlags(s.T(),
		[]string{
			"commit",
			"semver",
		},
		getVersionFlags(),
	)
}

func (s *CLIVersionHandlerTestSuite) Test_getVersionAction() {
	t := s.T()
	config := Config{}
	logger := InitLogger(&LoggerConfig{Name: "getVersionAction", Format: "raw", Level: "trace"})
	s.mockApp.Action = getVersionAction(&config, logger)
	if err := s.mockApp.Run([]string{"test-run-version"}); err == nil {
		assert.True(t, config.RunVersion)
	} else {
		panic(err)
	}
}
