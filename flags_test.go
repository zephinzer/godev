package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli"
)

type FlagsTestSuite struct {
	suite.Suite
}

func TestFlags(t *testing.T) {
	suite.Run(t, new(FlagsTestSuite))
}

func (s *FlagsTestSuite) Test_getFlagBuildOutput() {
	ensureFlag(s.T(), getFlagBuildOutput(), cli.StringFlag{}, `^output.*`)
}

func (s *FlagsTestSuite) Test_getFlagCommandArguments() {
	ensureFlag(s.T(), getFlagCommandArguments(), cli.StringFlag{}, `^args`)
}

func (s *FlagsTestSuite) Test_getFlagCommandsDelimiter() {
	ensureFlag(s.T(), getFlagCommandsDelimiter(), cli.StringFlag{}, `^exec-delim.*`)
}

func (s *FlagsTestSuite) Test_getFlagEnvVars() {
	ensureFlag(s.T(), getFlagEnvVars(), cli.StringSliceFlag{}, `^env.*`)
}

func (s *FlagsTestSuite) Test_getFlagExecGroups() {
	ensureFlag(s.T(), getFlagExecGroups(), cli.StringSliceFlag{}, `^exec.*`)
}

func (s *FlagsTestSuite) Test_getFlagFileExtensions() {
	ensureFlag(s.T(), getFlagFileExtensions(), cli.StringFlag{}, `^exts.*`)
}

func (s *FlagsTestSuite) Test_getFlagIgnoredNames() {
	ensureFlag(s.T(), getFlagIgnoredNames(), cli.StringFlag{}, `^ignore.*`)
}

func (s *FlagsTestSuite) Test_getFlagRate() {
	ensureFlag(s.T(), getFlagRate(), cli.DurationFlag{}, `^rate.*`)
}

func (s *FlagsTestSuite) Test_getFlagWatchDirectory() {
	ensureFlag(s.T(), getFlagWatchDirectory(), cli.StringFlag{}, `^watch.*`)
}

func (s *FlagsTestSuite) Test_getFlagWorkDirectory() {
	ensureFlag(s.T(), getFlagWorkDirectory(), cli.StringFlag{}, `^dir.*`)
}

func (s *FlagsTestSuite) Test_getFlagCommit() {
	ensureFlag(s.T(), getFlagCommit(), cli.BoolFlag{}, `^commit.*`)
}

func (s *FlagsTestSuite) Test_getFlagSemver() {
	ensureFlag(s.T(), getFlagSemver(), cli.BoolFlag{}, `^semver.*`)
}

func (s *FlagsTestSuite) Test_getFlagSilent() {
	ensureFlag(s.T(), getFlagSilent(), cli.BoolFlag{}, `^silent.*`)
}

func (s *FlagsTestSuite) Test_getFlagVerboseLogs() {
	ensureFlag(s.T(), getFlagVerboseLogs(), cli.BoolFlag{}, `^verbose.*`)
}

func (s *FlagsTestSuite) Test_getFlagSuperVerboseLogs() {
	ensureFlag(s.T(), getFlagSuperVerboseLogs(), cli.BoolFlag{}, `^vverbose.*`)
}
