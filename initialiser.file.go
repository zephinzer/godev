package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
)

type FileInitialiserConfig struct {
	Data     []byte
	Path     string
	Question string
}

func InitFileInitialiser(config *FileInitialiserConfig) *FileInitialiser {
	fi := &FileInitialiser{
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
	Data     []byte
	Path     string
	Question string
	handler  func() error
	logger   *Logger
}

func (fi FileInitialiser) Check() bool {
	return fileExists(fi.Path)
}

func (fi FileInitialiser) Confirm(reader *bufio.Reader) bool {
	return confirm(
		reader,
		Color("white", "godev> "+fi.Question),
		false,
		Color("bold", Color("red", InitialiserRetryText)),
	)
}

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
	} else {
		if file, err := os.Create(fi.Path); err != nil {
			return err
		} else if _, err = file.Write(fi.Data); err != nil {
			return err
		}
		return nil
	}
}
