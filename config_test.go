package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (s *ConfigTestSuite) Test_assignDefaultsRun() {
	t := s.T()
	c := &Config{
		BuildOutput:    "bin/app",
		WatchDirectory: "/some/path",
	}
	c.assignDefaults()
	assert.Equal(t, "info", c.LogLevel.String())
	assert.Equal(t, "/some/path/bin/app", c.BuildOutput)
	assert.Equal(t, false, c.RunView)
	assert.Equal(t, []string{"bin", "vendor"}, []string(c.IgnoredNames))
	assert.Equal(t, []string{"go", "Makefile"}, []string(c.FileExtensions))
	assert.Equal(t, "go mod vendor", c.ExecGroups[0])
	assert.Equal(t, "go build -o /some/path/bin/app", c.ExecGroups[1])
	assert.Equal(t, "/some/path/bin/app", c.ExecGroups[2])
}

func (s *ConfigTestSuite) Test_assignDefaultsTest() {
	t := s.T()
	c := &Config{
		BuildOutput:    "bin/app",
		WatchDirectory: "/some/path",
		RunTest:        true,
	}
	c.assignDefaults()
	assert.Equal(t, "info", c.LogLevel.String())
	assert.Equal(t, "/some/path/bin/app", c.BuildOutput)
	assert.Equal(t, false, c.RunView)
	assert.Equal(t, []string{"bin", "vendor"}, []string(c.IgnoredNames))
	assert.Equal(t, []string{"go", "Makefile"}, []string(c.FileExtensions))
	assert.Equal(t, "go mod vendor", c.ExecGroups[0])
	assert.Equal(t, "go build -o /some/path/bin/app", c.ExecGroups[1])
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
