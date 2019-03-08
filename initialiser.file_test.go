package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FileInitialiserTestSuite struct {
	suite.Suite
	expectedFileContent     string
	fileInitialiser         *FileInitialiser
	filePathThatExists      string
	filePathThatDoesntExist string
	logs                    bytes.Buffer
}

func TestFileInitialiserTestSuite(t *testing.T) {
	suite.Run(t, new(FileInitialiserTestSuite))
}

func (s *FileInitialiserTestSuite) SetupTest() {
	s.filePathThatExists = path.Join(getCurrentWorkingDirectory(), "/data/test-initialiser/file/exists")
	s.filePathThatDoesntExist = path.Join(getCurrentWorkingDirectory(), "/data/test-initialiser/file/non-existent")
	s.expectedFileContent = "FileInitialiserTestSuite"
	if fileExists(s.filePathThatDoesntExist) {
		if err := os.Remove(s.filePathThatDoesntExist); err != nil {
			panic(err)
		}
	}
	s.fileInitialiser = InitFileInitialiser(&FileInitialiserConfig{
		Data:     []byte(s.expectedFileContent),
		Question: "question",
	})
	s.fileInitialiser.logger.SetOutput(&s.logs)
}

func (s *FileInitialiserTestSuite) TestInitFileInitialiser_initialisesALogger() {
	assert.NotNil(s.T(), s.fileInitialiser.logger)
}

func (s *FileInitialiserTestSuite) TestCheck() {
	s.fileInitialiser.Path = s.filePathThatExists
	assert.True(s.T(), s.fileInitialiser.Check())
	s.fileInitialiser.Path = s.filePathThatDoesntExist
	assert.False(s.T(), s.fileInitialiser.Check())
}

func (s *FileInitialiserTestSuite) TestConfirm() {
	reader := bufio.NewReader(strings.NewReader("y\n"))
	assert.True(s.T(), s.fileInitialiser.Confirm(reader))
	reader = bufio.NewReader(strings.NewReader("n\n"))
	assert.False(s.T(), s.fileInitialiser.Confirm(reader))
}

func (s *FileInitialiserTestSuite) TestHandle_skip() {
	s.fileInitialiser.Path = s.filePathThatDoesntExist
	s.fileInitialiser.Handle(true)
	assert.Contains(
		s.T(),
		s.logs.String(),
		fmt.Sprintf(
			"skipping '%s'",
			path.Base(s.fileInitialiser.Path),
		),
	)
}

func (s *FileInitialiserTestSuite) TestHandle_customHandler() {
	s.fileInitialiser.Path = s.filePathThatDoesntExist
	s.fileInitialiser.handler = func() error {
		return errors.New("should be called")
	}
	assert.Equal(s.T(), "should be called", s.fileInitialiser.Handle().Error())
}

func (s *FileInitialiserTestSuite) TestHandle_defaultHandler() {
	s.fileInitialiser.Path = s.filePathThatDoesntExist
	s.fileInitialiser.Handle()
	assert.True(s.T(), fileExists(s.filePathThatDoesntExist))
	if contents, err := ioutil.ReadFile(s.filePathThatDoesntExist); err != nil {
		panic(err)
	} else {
		assert.Equal(s.T(), string(contents), s.expectedFileContent)
	}
	if err := os.Remove(s.filePathThatDoesntExist); err != nil {
		panic(err)
	}
}
