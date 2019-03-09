package main

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LoggerTestSuite struct {
	suite.Suite
	logger *Logger
	logs   bytes.Buffer
}

func TestLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}

func (s *LoggerTestSuite) SetupTest() {
	s.logger = InitLogger(&LoggerConfig{
		Name:   "LoggerTestSuite",
		Format: "production",
		Level:  "trace",
	})
	s.logger.SetOutput(&s.logs)
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

func (s *LoggerTestSuite) TestTrace() {
	s.logger.Trace("Trace")
	s.logger.Tracef("%s", "ecarT")
	assert.Contains(s.T(), s.logs.String(), "Trace")
	assert.Contains(s.T(), s.logs.String(), "ecarT")
}

func (s *LoggerTestSuite) TestDebug() {
	s.logger.Debug("Debug")
	s.logger.Debugf("%s", "gubeD")
	assert.Contains(s.T(), s.logs.String(), "Debug")
	assert.Contains(s.T(), s.logs.String(), "gubeD")
}

func (s *LoggerTestSuite) TestInfo() {
	s.logger.Info("Info")
	s.logger.Infof("%s", "ofnI")
	assert.Contains(s.T(), s.logs.String(), "Info")
	assert.Contains(s.T(), s.logs.String(), "ofnI")
}

func (s *LoggerTestSuite) TestWarn() {
	s.logger.Warn("Warn")
	s.logger.Warnf("%s", "nraW")
	assert.Contains(s.T(), s.logs.String(), "Warn")
	assert.Contains(s.T(), s.logs.String(), "nraW")
}

func (s *LoggerTestSuite) TestError() {
	s.logger.Error("Error")
	s.logger.Errorf("%s", "rorrE")
	assert.Contains(s.T(), s.logs.String(), "Error")
	assert.Contains(s.T(), s.logs.String(), "rorrE")
}

func (s *LoggerTestSuite) TestPanic() {
	defer func() {
		r := recover()
		assert.NotNil(s.T(), r)
	}()
	s.logger.Panic("Panic")
	s.logger.Panicf("%s", "cinaP")
	assert.Contains(s.T(), s.logs.String(), "Panic")
	assert.Contains(s.T(), s.logs.String(), "cinaP")
}
