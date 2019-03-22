package main

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
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
	flag := getFlagBuildOutput()
	assert.IsType(s.T(), cli.StringFlag{}, flag)
	stringFlag := cli.StringFlag(flag.(cli.StringFlag))
	assert.Regexp(s.T(), regexp.MustCompile(`^output.*`), stringFlag.Name)
}

func (s *FlagsTestSuite) Test_getFlagCommandsDelimiter() {
	flag := getFlagCommandsDelimiter()
	assert.IsType(s.T(), cli.StringFlag{}, flag)
	stringFlag := cli.StringFlag(flag.(cli.StringFlag))
	assert.Regexp(s.T(), regexp.MustCompile(`^exec-delim.*`), stringFlag.Name)
}

func (s *FlagsTestSuite) Test_getFlagEnvVars() {
	flag := getFlagEnvVars()
	assert.IsType(s.T(), cli.StringSliceFlag{}, flag)
	stringSliceFlag := cli.StringSliceFlag(flag.(cli.StringSliceFlag))
	assert.Regexp(s.T(), regexp.MustCompile(`^env.*`), stringSliceFlag.Name)
}

func (s *FlagsTestSuite) Test_getFlagExecGroups() {
	flag := getFlagExecGroups()
	assert.IsType(s.T(), cli.StringSliceFlag{}, flag)
	stringSliceFlag := cli.StringSliceFlag(flag.(cli.StringSliceFlag))
	assert.Regexp(s.T(), regexp.MustCompile(`^exec.*`), stringSliceFlag.Name)
}

func (s *FlagsTestSuite) Test_getFlagFileExtensions() {
	flag := getFlagFileExtensions()
	assert.IsType(s.T(), cli.StringFlag{}, flag)
	stringFlag := cli.StringFlag(flag.(cli.StringFlag))
	assert.Regexp(s.T(), regexp.MustCompile(`^exts.*`), stringFlag.Name)
}

func (s *FlagsTestSuite) Test_getFlagIgnoredNames() {
	flag := getFlagIgnoredNames()
	assert.IsType(s.T(), cli.StringFlag{}, flag)
	stringFlag := cli.StringFlag(flag.(cli.StringFlag))
	assert.Regexp(s.T(), regexp.MustCompile(`^ignore.*`), stringFlag.Name)
}

func (s *FlagsTestSuite) Test_getFlagRate() {
	flag := getFlagRate()
	assert.IsType(s.T(), cli.DurationFlag{}, flag)
	durationFlag := cli.DurationFlag(flag.(cli.DurationFlag))
	assert.Regexp(s.T(), regexp.MustCompile(`^rate.*`), durationFlag.Name)
}

func (s *FlagsTestSuite) Test_getFlagWatchDirectory() {
	flag := getFlagWatchDirectory()
	assert.IsType(s.T(), cli.StringFlag{}, flag)
	stringFlag := cli.StringFlag(flag.(cli.StringFlag))
	assert.Regexp(s.T(), regexp.MustCompile(`^watch.*`), stringFlag.Name)
}

func (s *FlagsTestSuite) Test_getFlagWorkDirectory() {
	flag := getFlagWorkDirectory()
	assert.IsType(s.T(), cli.StringFlag{}, flag)
	stringFlag := cli.StringFlag(flag.(cli.StringFlag))
	assert.Regexp(s.T(), regexp.MustCompile(`^dir.*`), stringFlag.Name)
}

func (s *FlagsTestSuite) Test_getFlagCommit() {
	flag := getFlagCommit()
	assert.IsType(s.T(), cli.BoolFlag{}, flag)
	boolFlag := cli.BoolFlag(flag.(cli.BoolFlag))
	assert.Regexp(s.T(), regexp.MustCompile(`^commit.*`), boolFlag.Name)
}

func (s *FlagsTestSuite) Test_getFlagSemver() {
	flag := getFlagSemver()
	assert.IsType(s.T(), cli.BoolFlag{}, flag)
	boolFlag := cli.BoolFlag(flag.(cli.BoolFlag))
	assert.Regexp(s.T(), regexp.MustCompile(`^semver.*`), boolFlag.Name)
}

func (s *FlagsTestSuite) Test_getFlagSilent() {
	flag := getFlagSilent()
	assert.IsType(s.T(), cli.BoolFlag{}, flag)
	boolFlag := cli.BoolFlag(flag.(cli.BoolFlag))
	assert.Regexp(s.T(), regexp.MustCompile(`^silent.*`), boolFlag.Name)
}

func (s *FlagsTestSuite) Test_getFlagVerboseLogs() {
	flag := getFlagVerboseLogs()
	assert.IsType(s.T(), cli.BoolFlag{}, flag)
	boolFlag := cli.BoolFlag(flag.(cli.BoolFlag))
	assert.Regexp(s.T(), regexp.MustCompile(`^verbose.*`), boolFlag.Name)
}

func (s *FlagsTestSuite) Test_getFlagSuperVerboseLogs() {
	flag := getFlagSuperVerboseLogs()
	assert.IsType(s.T(), cli.BoolFlag{}, flag)
	boolFlag := cli.BoolFlag(flag.(cli.BoolFlag))
	assert.Regexp(s.T(), regexp.MustCompile(`^vverbose.*`), boolFlag.Name)
}
