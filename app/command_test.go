package main

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CommandTestSuite struct {
	suite.Suite
}

func TestCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CommandTestSuite))
}

func (s *CommandTestSuite) TestIsValidFromRegisteredPath() {
	// we are running using `go` so there's no reason why it should be unavailable
	expectedApplication := "go"
	testCommand := InitCommand(&CommandConfig{
		Application: expectedApplication,
		Arguments:   []string{},
	})
	err := testCommand.IsValid()
	assert.Nil(s.T(), err)
}

func (s *CommandTestSuite) TestIsValidFromAbsolutePathNoPermissions() {
	cwd := getCurrentWorkingDirectory()
	expectedApplication := path.Join(cwd, "/data/test-exec/nonexec.sh")
	testCommand := InitCommand(&CommandConfig{
		Application: expectedApplication,
		Arguments:   []string{},
	})
	err := testCommand.IsValid()
	assert.NotNil(s.T(), err)
	assert.Contains(s.T(), err.Error(), "you don't have permissions to execute")
}

func (s *CommandTestSuite) TestIsValidFromAbsolutePathWithPermissions() {
	cwd := getCurrentWorkingDirectory()
	expectedApplication := path.Join(cwd, "/data/test-exec/exec.sh")
	testCommand := InitCommand(&CommandConfig{
		Application: expectedApplication,
		Arguments:   []string{},
	})
	err := testCommand.IsValid()
	assert.Nil(s.T(), err)
}
