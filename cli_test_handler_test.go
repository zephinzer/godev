package main

import (
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli"
)

type CLITestHandlerTestSuite struct {
	suite.Suite
	mockApp     *cli.App
	mockContext *cli.Context
}

func TestCLITestHandler(t *testing.T) {
	suite.Run(t, new(CLITestHandlerTestSuite))
}

func (s *CLITestHandlerTestSuite) SetupTest() {
	s.mockApp = cli.NewApp()
	s.mockApp.Flags = getTestFlags()
	s.mockContext = cli.NewContext(s.mockApp, nil, nil)
}

func (s *CLITestHandlerTestSuite) Test_getTestCommand() {
	config := Config{}
	command := getTestCommand(&config)
	ensureCLICommand(s.T(), command, "test", "t", getTestFlags())
}

func (s *CLITestHandlerTestSuite) Test_getTestFlags() {
	ensureCLIFlags(s.T(),
		[]string{
			"dir",
			"env",
			"exec-delim",
			"exts",
			"ignore",
			"output",
			"rate",
			"silent",
			"verbose",
			"vverbose",
			"watch",
		},
		getTestFlags(),
	)
}

func (s *CLITestHandlerTestSuite) Test_getTestAction() {
	t := s.T()
	config := Config{}
	s.mockApp.Action = getTestAction(&config)
	if err := s.mockApp.Run([]string{"test-run-test"}); err == nil {
		pathToBinary := path.Join(getCurrentWorkingDirectory(), "/bin/app")
		assert.True(t, config.RunTest)
		assert.Equal(t, pathToBinary, config.BuildOutput)
		assert.Equal(t, ",", config.CommandsDelimiter)
		assert.Equal(t, []string{}, []string(config.EnvVars))
		assert.Equal(t, []string{"go mod vendor", "go build -o " + pathToBinary, "go test ./... -coverprofile c.out"}, []string(config.ExecGroups))
		assert.Equal(t, []string{"go", "Makefile"}, []string(config.FileExtensions))
		assert.Equal(t, []string{"bin", "vendor"}, []string(config.IgnoredNames))
		assert.Equal(t, 2*time.Second, config.Rate)
		assert.Equal(t, getCurrentWorkingDirectory(), config.WatchDirectory)
		assert.Equal(t, getCurrentWorkingDirectory(), config.WorkDirectory)
		assert.Equal(t, "info", string(config.LogLevel))
	} else {
		panic(err)
	}
}
