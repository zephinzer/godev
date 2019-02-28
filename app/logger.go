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

type rawFormat struct{}

func (f *rawFormat) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}

type productionFormat struct{}

func (f *productionFormat) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("Jan02/15:04")
	data := entry.Data
	var moduleLabel string
	moduleName := data["module"]
	if moduleName != nil {
		moduleLabel = fmt.Sprintf("%v", moduleName)
	}
	var submoduleLabel string
	submoduleName := data["submodule"]
	if submoduleName != nil {
		submoduleLabel = fmt.Sprintf("/%v", submoduleName)
	}
	message := fmt.Sprintf("|%v| [%v%v] %s", timestamp, moduleLabel, submoduleLabel, entry.Message)
	if entry.Level > logrus.InfoLevel {
		var otherKeys string
		for key, value := range data {
			if key != "module" && key != "submodule" {
				otherKeys = fmt.Sprintf("%s\n  %s: %v", otherKeys, key, value)
			}
		}
		message = fmt.Sprintf("%s%s", message, otherKeys)
	}
	var log []byte
	switch entry.Level {
	case logrus.TraceLevel:
		log = []byte(Color("gray", fmt.Sprintf("%s\n", message)))
	case logrus.DebugLevel:
		log = []byte(Color("bold", Color("gray", fmt.Sprintf("%s\n", message))))
	case logrus.InfoLevel:
		log = []byte(Color("green", fmt.Sprintf("%s\n", message)))
	case logrus.WarnLevel:
		log = []byte(Color("yellow", fmt.Sprintf("%s\n", message)))
	case logrus.ErrorLevel:
		log = []byte(Color("red", fmt.Sprintf("%s\n", message)))
	case logrus.PanicLevel:
		log = []byte(Color("bold", Color("red", fmt.Sprintf("%s\n", message))))
	default:
		log = []byte(Color("default", fmt.Sprintf("%s\n", message)))
	}
	return log, nil
}

func (lf *LogFormat) Get() logrus.Formatter {
	switch *lf {
	case "json":
		return &logrus.JSONFormatter{}
	case "production":
		return new(productionFormat)
	case "raw":
		return new(rawFormat)
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
	Name             string
	Format           LogFormat
	Level            LogLevel
	AdditionalFields *map[string]interface{}
}

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

type Logger struct {
	config      *LoggerConfig
	instance    *logrus.Entry
	instanceRaw *logrus.Logger
}

// SetOutput exists for characterisation testing
func (l *Logger) SetOutput(buffer *bytes.Buffer) {
	l.instanceRaw.SetOutput(buffer)
}

// WithFields - enables adding of more fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	l.instanceRaw.WithFields(logrus.Fields(fields))
	return l
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
