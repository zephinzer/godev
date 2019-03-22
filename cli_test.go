package main

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CLITestSuite struct {
	suite.Suite
}

func TestCLI(t *testing.T) {
	suite.Run(t, new(CLITestSuite))
}

func (s *CLITestSuite) Test_initCLI() {
	cli := initCLI()
	assert.NotNil(s.T(), cli.config)
	assert.NotNil(s.T(), cli.logger)
	assert.NotNil(s.T(), cli.instance)
	assert.NotNil(s.T(), cli.instance.Name)
	assert.NotNil(s.T(), cli.instance.Usage)
	assert.NotNil(s.T(), cli.instance.Description)
}

func (s *CLITestSuite) TestStart_provisionsDefault() {
	cli := initCLI()
	cli.Start([]string{"godev"}, func(config *Config) {
		assert.NotNil(s.T(), config)
		assert.Equal(s.T(), true, config.RunDefault)
	})
}

func (s *CLITestSuite) TestStart_provisionsInit() {
	cli := initCLI()
	cli.Start([]string{"godev", "init"}, func(config *Config) {
		assert.NotNil(s.T(), config)
		assert.Equal(s.T(), true, config.RunInit)
	})
}

func (s *CLITestSuite) TestStart_provisionsTest() {
	cli := initCLI()
	cli.Start([]string{"godev", "test"}, func(config *Config) {
		assert.NotNil(s.T(), config)
		assert.Equal(s.T(), true, config.RunTest)
	})
}

func (s *CLITestSuite) TestStart_provisionsVersion() {
	cli := initCLI()
	var logs bytes.Buffer
	cli.rawLogger.SetOutput(&logs)
	cli.Start([]string{"godev", "version"}, func(config *Config) {
		assert.NotNil(s.T(), config)
		assert.Equal(s.T(), true, config.RunVersion)
		assert.Regexp(s.T(), regexp.MustCompile(`^godev [\d]*\.[\d]*\.[\d]*\-[a-f0-9]{7}`), logs.String())
	})
}

func (s *CLITestSuite) TestStart_provisionsView() {
	cli := initCLI()
	var logs bytes.Buffer
	cli.rawLogger.SetOutput(&logs)
	cli.Start([]string{"godev", "view", "dockerfile"}, func(config *Config) {
		assert.NotNil(s.T(), config)
		assert.Equal(s.T(), true, config.RunView)
		assert.Contains(s.T(), logs.String(), DataDockerfile)
	})
	cli.Start([]string{"godev", "view", "makefile"}, func(config *Config) {
		assert.NotNil(s.T(), config)
		assert.Equal(s.T(), true, config.RunView)
		assert.Contains(s.T(), logs.String(), DataMakefile)
	})
}
