package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

// FileInitialiserConfig holds the configuration for initialising the files
type FileInitialiserConfig struct {
	Data     []byte
	Path     string
	Question string
}

// InitFileInitialiser is the method for creating a file initialiser
func InitFileInitialiser(config *FileInitialiserConfig) *FileInitialiser {
	fi := &FileInitialiser{
		Key:      strings.ToLower(path.Base(config.Path)),
		Data:     config.Data,
		Path:     config.Path,
		Question: config.Question,
		logger: InitLogger(&LoggerConfig{
			Format: "raw",
		}),
	}
	return fi
}

// FileInitialiser is used for checking if certain files
// exist at :Path and questioning the user if they'd like to
// seed it
type FileInitialiser struct {
	Key      string
	Data     []byte
	Path     string
	Question string
	handler  func() error
	logger   *Logger
}

// Check verifies if the Question should be popped
func (fi FileInitialiser) Check() bool {
	return fileExists(fi.Path)
}

// Confirm seeks advice from the user whether we should proceed
func (fi FileInitialiser) Confirm(reader *bufio.Reader) bool {
	return confirm(
		reader,
		Color("white", "godev> "+fi.Question),
		false,
		Color("bold", Color("red", initialiserRetryText)),
	)
}

// GetKey returns the key of this file initialiser
func (fi *FileInitialiser) GetKey() string {
	return fi.Key
}

// Handle initialises the file if it doesn't exist or if :skip is indiciated
func (fi FileInitialiser) Handle(skip ...bool) error {
	if len(skip) > 0 && skip[0] {
		fi.logger.Info(
			Color("gray",
				fmt.Sprintf("godev> skipping '%s' - already exists", path.Base(fi.Path)),
			),
		)
		return nil
	}
	if fi.handler != nil {
		return fi.handler()
	}
	if file, err := os.Create(fi.Path); err != nil {
		return err
	} else if _, err = file.Write(fi.Data); err != nil {
		return err
	}
	return nil
}
