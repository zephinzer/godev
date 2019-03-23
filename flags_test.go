package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type FlagsTestSuite struct {
	suite.Suite
}

func TestFlags(t *testing.T) {
	suite.Run(t, new(FlagsTestSuite))
}

func (s *FlagsTestSuite) Test_getFlagBuildOutput() {
	ensureStringFlag(s.T(), getFlagBuildOutput(), `^output.*`)
}

func (s *FlagsTestSuite) Test_getFlagCommandsDelimiter() {
	ensureStringFlag(s.T(), getFlagCommandsDelimiter(), `^exec-delim.*`)
}

func (s *FlagsTestSuite) Test_getFlagEnvVars() {
	ensureStringSliceFlag(s.T(), getFlagEnvVars(), `^env.*`)
}

func (s *FlagsTestSuite) Test_getFlagExecGroups() {
	ensureStringSliceFlag(s.T(), getFlagExecGroups(), `^exec.*`)
}

func (s *FlagsTestSuite) Test_getFlagFileExtensions() {
	ensureStringFlag(s.T(), getFlagFileExtensions(), `^exts.*`)
}

func (s *FlagsTestSuite) Test_getFlagIgnoredNames() {
	ensureStringFlag(s.T(), getFlagIgnoredNames(), `^ignore.*`)
}

func (s *FlagsTestSuite) Test_getFlagRate() {
	ensureDurationFlag(s.T(), getFlagRate(), `^rate.*`)
}

func (s *FlagsTestSuite) Test_getFlagWatchDirectory() {
	ensureStringFlag(s.T(), getFlagWatchDirectory(), `^watch.*`)
}

func (s *FlagsTestSuite) Test_getFlagWorkDirectory() {
	ensureStringFlag(s.T(), getFlagWorkDirectory(), `^dir.*`)
}

func (s *FlagsTestSuite) Test_getFlagCommit() {
	ensureBoolFlag(s.T(), getFlagCommit(), `^commit.*`)
}

func (s *FlagsTestSuite) Test_getFlagSemver() {
	ensureBoolFlag(s.T(), getFlagSemver(), `^semver.*`)
}

func (s *FlagsTestSuite) Test_getFlagSilent() {
	ensureBoolFlag(s.T(), getFlagSilent(), `^silent.*`)
}

func (s *FlagsTestSuite) Test_getFlagVerboseLogs() {
	ensureBoolFlag(s.T(), getFlagVerboseLogs(), `^verbose.*`)
}

func (s *FlagsTestSuite) Test_getFlagSuperVerboseLogs() {
	ensureBoolFlag(s.T(), getFlagSuperVerboseLogs(), `^vverbose.*`)
}
