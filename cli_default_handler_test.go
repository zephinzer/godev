package main

import (
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli"
)

type CLIDefaultHandlerTestSuite struct {
	suite.Suite
	mockApp     *cli.App
	mockContext *cli.Context
}

func TestCLIDefaultHandler(t *testing.T) {
	suite.Run(t, new(CLIDefaultHandlerTestSuite))
}

func (s *CLIDefaultHandlerTestSuite) SetupTest() {
	s.mockApp = cli.NewApp()
	s.mockApp.Flags = getDefaultFlags()
	s.mockContext = cli.NewContext(s.mockApp, nil, nil)
}

func (s *CLIDefaultHandlerTestSuite) Test_getDefaultFlags() {
	ensureCLIFlags(s.T(),
		[]string{
			"dir",
			"env",
			"exec-delim",
			"exec",
			"exts",
			"ignore",
			"output",
			"rate",
			"silent",
			"verbose",
			"vverbose",
			"watch",
		},
		getDefaultFlags(),
	)
}

func (s *CLIDefaultHandlerTestSuite) Test_getDefaultAction() {
	t := s.T()
	config := Config{}
	s.mockApp.Action = getDefaultAction(&config)
	if err := s.mockApp.Run([]string{"test-run"}); err == nil {
		pathToBinary := path.Join(getCurrentWorkingDirectory(), "/bin/app")
		assert.True(t, config.RunDefault)
		assert.Equal(t, pathToBinary, config.BuildOutput)
		assert.Equal(t, ",", config.CommandsDelimiter)
		assert.Equal(t, []string{}, []string(config.EnvVars))
		assert.Equal(t, []string{"go mod vendor", "go build -o " + pathToBinary, pathToBinary}, []string(config.ExecGroups))
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
