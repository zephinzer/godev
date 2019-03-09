package main

import (
	"path"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (s *ConfigTestSuite) TestInitConfig() {
	t := s.T()
	config := InitConfig()
	assert.Equal(t, path.Join(getCurrentWorkingDirectory(), DefaultBuildOutput), config.BuildOutput)
	assert.Equal(t, DefaultCommandsDelimiter, config.CommandsDelimiter)
	assert.Len(t, config.ExecGroups, 3)
	assert.Contains(t, config.ExecGroups[0], "go mod")
	assert.Contains(t, config.ExecGroups[1], "go build")
	assert.Contains(t, config.ExecGroups[2], DefaultBuildOutput)
	assert.Equal(t, strings.Split(DefaultFileExtensions, ","), []string(config.FileExtensions))
	assert.Equal(t, strings.Split(DefaultIgnoredNames, ","), []string(config.IgnoredNames))
	assert.Equal(t, false, config.LogSilent)
	assert.Equal(t, false, config.LogVerbose)
	assert.Equal(t, false, config.LogSuperVerbose)
	assert.Equal(t, 2*time.Second, config.Rate)
	assert.Equal(t, false, config.RunInit)
	assert.Equal(t, false, config.RunTest)
	assert.Equal(t, false, config.RunVersion)
	assert.Equal(t, false, config.RunView)
	assert.Equal(t, "", config.View)
	assert.Equal(t, getCurrentWorkingDirectory(), config.WatchDirectory)
	assert.Equal(t, getCurrentWorkingDirectory(), config.WorkDirectory)
	assert.Equal(t, DefaultLogLevel, string(config.LogLevel))
}

func (s *ConfigTestSuite) Test_assignDefaultsRun() {
	t := s.T()
	c := &Config{
		BuildOutput:    "bin/app",
		WatchDirectory: "/some/path/to/watch",
		WorkDirectory:  "/some/path/to/work",
	}
	c.assignDefaults()
	assert.Equal(t, "info", c.LogLevel.String())
	assert.Equal(t, "/some/path/to/work/bin/app", c.BuildOutput)
	assert.Equal(t, false, c.RunView)
	assert.Equal(t, []string{"bin", "vendor"}, []string(c.IgnoredNames))
	assert.Equal(t, []string{"go", "Makefile"}, []string(c.FileExtensions))
	assert.Equal(t, "go mod vendor", c.ExecGroups[0])
	assert.Equal(t, "go build -o /some/path/to/work/bin/app", c.ExecGroups[1])
	assert.Equal(t, "/some/path/to/work/bin/app", c.ExecGroups[2])
}

func (s *ConfigTestSuite) Test_assignDefaultsTest() {
	t := s.T()
	c := &Config{
		BuildOutput:    "bin/app",
		WatchDirectory: "/some/path/to/watch",
		WorkDirectory:  "/some/path/to/work",
		RunTest:        true,
	}
	c.assignDefaults()
	assert.Equal(t, "info", c.LogLevel.String())
	assert.Equal(t, "/some/path/to/work/bin/app", c.BuildOutput)
	assert.Equal(t, false, c.RunView)
	assert.Equal(t, []string{"bin", "vendor"}, []string(c.IgnoredNames))
	assert.Equal(t, []string{"go", "Makefile"}, []string(c.FileExtensions))
	assert.Equal(t, "go mod vendor", c.ExecGroups[0])
	assert.Equal(t, "go build -o /some/path/to/work/bin/app", c.ExecGroups[1])
	assert.Equal(t, "go test ./... -coverprofile c.out", c.ExecGroups[2])
}

func (s *ConfigTestSuite) Test_interpretLogLevel() {
	c := &Config{LogVerbose: true}
	c.interpretLogLevel()
	assert.Equal(s.T(), "debug", c.LogLevel.String())
	c = &Config{LogSuperVerbose: true}
	c.interpretLogLevel()
	assert.Equal(s.T(), "trace", c.LogLevel.String())
}
