package main

import (
	"bytes"

	"github.com/sirupsen/logrus"
)

// LoggerConfig is used for initialising the logger
type LoggerConfig struct {
	Name             string
	Format           LogFormat
	Level            LogLevel
	AdditionalFields *map[string]interface{}
}

// InitLogger is used for setting up a new logger for a component
func InitLogger(config *LoggerConfig) *Logger {
	log := logrus.New()
	log.SetFormatter(config.Format.Get())
	log.SetLevel(config.Level.Get())
	fields := logrus.Fields{
		"module": config.Name,
	}
	if config.AdditionalFields != nil {
		for key, value := range *config.AdditionalFields {
			fields[key] = value
		}
	}
	logger := &Logger{
		config:      config,
		instanceRaw: log,
		instance:    log.WithFields(fields),
	}
	return logger
}

// Logger represents the logger instance created on InitLogger
type Logger struct {
	config      *LoggerConfig
	instance    *logrus.Entry
	instanceRaw *logrus.Logger
}

// SetOutput exists for characterisation testing
func (l *Logger) SetOutput(buffer *bytes.Buffer) {
	l.instanceRaw.SetOutput(buffer)
}

// Trace logs at the trace level
func (l *Logger) Trace(log ...interface{}) {
	l.instance.Trace(log...)
}

// Tracef logs at the trace level with formatting
func (l *Logger) Tracef(format string, log ...interface{}) {
	l.instance.Tracef(format, log...)
}

// Debug logs at the debug level
func (l *Logger) Debug(log ...interface{}) {
	l.instance.Debug(log...)
}

// Debugf logs at the debug level with formatting
func (l *Logger) Debugf(format string, log ...interface{}) {
	l.instance.Debugf(format, log...)
}

// Info logs at the info level
func (l *Logger) Info(log ...interface{}) {
	l.instance.Info(log...)
}

// Infof logs at the info level with formatting
func (l *Logger) Infof(format string, log ...interface{}) {
	l.instance.Infof(format, log...)
}

// Warn logs at the warn level
func (l *Logger) Warn(log ...interface{}) {
	l.instance.Warn(log...)
}

// Warnf logs at the warn level with formatting
func (l *Logger) Warnf(format string, log ...interface{}) {
	l.instance.Warnf(format, log...)
}

// Error logs at the error level
func (l *Logger) Error(log ...interface{}) {
	l.instance.Error(log...)
}

// Errorf logs at the error level with formatting
func (l *Logger) Errorf(format string, log ...interface{}) {
	l.instance.Errorf(format, log...)
}

// Fatal logs at the fatal level
func (l *Logger) Fatal(log ...interface{}) {
	l.instance.Fatal(log...)
}

// Fatalf logs at the fatal level with formatting
func (l *Logger) Fatalf(format string, log ...interface{}) {
	l.instance.Fatalf(format, log...)
}

// Panic logs at the panic level
func (l *Logger) Panic(log ...interface{}) {
	l.instance.Panic(log...)
}

// Panicf logs at the panic level with formatting
func (l *Logger) Panicf(format string, log ...interface{}) {
	l.instance.Panicf(format, log...)
}

// LogLevel is a string represent of the log level
type LogLevel string

// String implements the string return type for LogLevel
func (ll *LogLevel) String() string {
	return string(*ll)
}

// Get retrieves the LogLevel for logrus to use
func (ll *LogLevel) Get() logrus.Level {
	switch *ll {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.TraceLevel
	}
}
