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
	verifyStart func(*testing.T, *bool)
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
	ensureCLIStartSetsRunFlag(s.T(), []string{"godev"}, "RunDefault")
}

func (s *CLITestSuite) TestStart_provisionsInit() {
	ensureCLIStartSetsRunFlag(s.T(), []string{"godev", "init"}, "RunInit")
}

func (s *CLITestSuite) TestStart_provisionsTest() {
	ensureCLIStartSetsRunFlag(s.T(), []string{"godev", "test"}, "RunTest")
}

func (s *CLITestSuite) TestStart_provisionsVersion() {
	ensureCLIStartSetsRunFlag(s.T(), []string{"godev", "version"}, "RunVersion", func(logs bytes.Buffer) {
		assert.Regexp(
			s.T(),
			regexp.MustCompile(`^godev [\d]*\.[\d]*\.[\d]*\-[a-f0-9]{7}`),
			logs.String(),
		)
	})
}

func (s *CLITestSuite) TestStart_provisionsView() {
	ensureCLIStartSetsRunFlag(s.T(), []string{"godev", "view", "dockerfile"}, "RunView", func(logs bytes.Buffer) {
		assert.Contains(s.T(), logs.String(), DataDockerfile)
	})
	ensureCLIStartSetsRunFlag(s.T(), []string{"godev", "view", "makefile"}, "RunView", func(logs bytes.Buffer) {
		assert.Contains(s.T(), logs.String(), DataMakefile)
	})
}
