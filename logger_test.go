package main

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LoggerTestSuite struct {
	suite.Suite
}

func TestLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}

func (s *LoggerTestSuite) TestLogLevelGet() {
	logLevel := LogLevel("trace")
	assert.Equal(s.T(), logLevel.Get(), logrus.TraceLevel)
	logLevel = LogLevel("debug")
	assert.Equal(s.T(), logLevel.Get(), logrus.DebugLevel)
	logLevel = LogLevel("info")
	assert.Equal(s.T(), logLevel.Get(), logrus.InfoLevel)
	logLevel = LogLevel("warn")
	assert.Equal(s.T(), logLevel.Get(), logrus.WarnLevel)
	logLevel = LogLevel("error")
	assert.Equal(s.T(), logLevel.Get(), logrus.ErrorLevel)
	logLevel = LogLevel("fatal")
	assert.Equal(s.T(), logLevel.Get(), logrus.FatalLevel)
	logLevel = LogLevel("something else")
	assert.Equal(s.T(), logLevel.Get(), logrus.TraceLevel)
}
