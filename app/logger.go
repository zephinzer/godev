package main

import (
	"bytes"
	"fmt"

	"github.com/sirupsen/logrus"
)

type LogFormat string

func (lf *LogFormat) String() string {
	if *lf == "json" {
		return "json"
	}
	return "text"
}

type productionFormat struct{}

func (f *productionFormat) Format(entry *logrus.Entry) ([]byte, error) {
	switch entry.Level {
	case logrus.TraceLevel:
		return []byte(Color(CDarkGray, fmt.Sprintf("'| [%v] %s\n", entry.Data, entry.Message))), nil
	case logrus.DebugLevel:
		return []byte(Color(CLightGray, fmt.Sprintf(".| [%v] %s\n", entry.Data, entry.Message))), nil
	case logrus.InfoLevel:
		return []byte(Color(CGreen, fmt.Sprintf(">| [%v] %s\n", entry.Data, entry.Message))), nil
	case logrus.WarnLevel:
		return []byte(Color(CYellow, fmt.Sprintf("!| [%v] %s\n", entry.Data, entry.Message))), nil
	case logrus.ErrorLevel:
		return []byte(Color(CRed, fmt.Sprintf("x| [%v] %s\n", entry.Data, entry.Message))), nil
	case logrus.PanicLevel:
		return []byte(Color(CRed, fmt.Sprintf("X| [%v] %s\n", entry.Data, entry.Message))), nil
	default:
		return []byte(Color(CDefault, fmt.Sprintf("-| [%v] %s\n", entry.Data, entry.Message))), nil
	}
}

func (lf *LogFormat) Get() logrus.Formatter {
	switch *lf {
	case "json":
		return &logrus.JSONFormatter{}
	case "production":
		return new(productionFormat)
	default:
		return &logrus.TextFormatter{
			ForceColors: true,
		}
	}
}

type LogLevel string

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
	default:
		return logrus.TraceLevel
	}
}

type LoggerConfig struct {
	Name   string
	Format LogFormat
	Level  LogLevel
}

func InitLogger(config *LoggerConfig) *Logger {
	log := logrus.New()
	log.SetFormatter(config.Format.Get())
	log.SetLevel(config.Level.Get())
	logger := &Logger{
		config:      config,
		instanceRaw: log,
		instance: log.WithFields(logrus.Fields{
			"module": config.Name,
		}),
	}
	return logger
}

type Logger struct {
	config      *LoggerConfig
	instance    *logrus.Entry
	instanceRaw *logrus.Logger
}

// SetOutput exists for characterisation testing
func (l *Logger) SetOutput(buffer *bytes.Buffer) {
	l.instanceRaw.SetOutput(buffer)
}

func (l *Logger) Trace(log ...interface{}) {
	l.instance.Trace(log...)
}

func (l *Logger) Tracef(format string, log ...interface{}) {
	l.instance.Tracef(format, log...)
}
func (l *Logger) Debug(log ...interface{}) {
	l.instance.Debug(log...)
}

func (l *Logger) Debugf(format string, log ...interface{}) {
	l.instance.Debugf(format, log...)
}

func (l *Logger) Info(log ...interface{}) {
	l.instance.Info(log...)
}

func (l *Logger) Infof(format string, log ...interface{}) {
	l.instance.Infof(format, log...)
}

func (l *Logger) Warn(log ...interface{}) {
	l.instance.Warn(log...)
}

func (l *Logger) Warnf(format string, log ...interface{}) {
	l.instance.Warnf(format, log...)
}

func (l *Logger) Error(log ...interface{}) {
	l.instance.Error(log...)
}

func (l *Logger) Errorf(format string, log ...interface{}) {
	l.instance.Errorf(format, log...)
}

func (l *Logger) Fatal(log ...interface{}) {
	l.instance.Fatal(log...)
}

func (l *Logger) Fatalf(format string, log ...interface{}) {
	l.instance.Fatalf(format, log...)
}

func (l *Logger) Panic(log ...interface{}) {
	l.instance.Panic(log...)
}

func (l *Logger) Panicf(format string, log ...interface{}) {
	l.instance.Panicf(format, log...)
}
