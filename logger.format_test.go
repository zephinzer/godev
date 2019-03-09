package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRawFormat(t *testing.T) {
	var logs bytes.Buffer
	logger := InitLogger(&LoggerConfig{
		Name:   "LoggerTestSuite",
		Format: "raw",
		Level:  "trace",
	})
	logger.SetOutput(&logs)
	logger.Debug("hi")
	logger.Trace("hi")
	logger.Info("hi")
	logger.Warn("hi")
	logger.Error("hi")
	assert.Equal(t, logs.String(), "hi\nhi\nhi\nhi\nhi\n")
	assert.NotContains(t, logs.String(), ColorStub)
}

func TestProductionFormat(t *testing.T) {
	var logs bytes.Buffer
	logger := InitLogger(&LoggerConfig{
		Name:   "LoggerTestSuite",
		Format: "production",
		Level:  "trace",
	})
	logger.SetOutput(&logs)
	logger.Debug("debug")
	assert.Contains(t, logs.String(), fmt.Sprintf("%s%vm", ColorStub, Palette["gray"]))
	logger.Info("info")
	assert.Contains(t, logs.String(), fmt.Sprintf("%s%vm", ColorStub, Palette["green"]))
	logger.Warn("warn")
	assert.Contains(t, logs.String(), fmt.Sprintf("%s%vm", ColorStub, Palette["yellow"]))
	logger.Error("error")
	assert.Contains(t, logs.String(), fmt.Sprintf("%s%vm", ColorStub, Palette["red"]))
	assert.Contains(t, logs.String(), "[LoggerTestSuite]")

}
