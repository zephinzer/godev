package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// LogFormat represents the possible string types that
// the log format can accept
type LogFormat string

// String returns a string representation of the format
func (lf *LogFormat) String() string {
	return string(*lf)
}

// Get returns a formatter which logrus can use
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

type rawFormat struct{}

func (f *rawFormat) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}

type productionFormat struct{}

var productionFormatColorMap = map[logrus.Level]string{
	logrus.TraceLevel: "gray",
	logrus.DebugLevel: "gray",
	logrus.InfoLevel:  "green",
	logrus.WarnLevel:  "yellow",
	logrus.ErrorLevel: "red",
	logrus.FatalLevel: "lred",
}

func (f *productionFormat) Format(entry *logrus.Entry) ([]byte, error) {
	var moduleLabel string
	var submoduleLabel string
	timestamp := entry.Time.Format("Jan02/15:04")
	data := entry.Data
	if data["module"] != nil {
		moduleLabel = fmt.Sprintf("%v", data["module"])
	}
	if data["submodule"] != nil {
		submoduleLabel = fmt.Sprintf("/%v", data["submodule"])
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
	color := productionFormatColorMap[entry.Level]
	log := []byte(Color(color, fmt.Sprintf("%s\n", message)))
	return log, nil
}
