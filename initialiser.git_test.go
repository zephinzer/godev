package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GitInitialiserTestSuite struct {
	suite.Suite
	gitInitialiser *GitInitialiser
	pathWithGit    string
	pathWithoutGit string
	logs           bytes.Buffer
}

func TestGitInitialiserTestSuite(t *testing.T) {
	suite.Run(t, new(GitInitialiserTestSuite))
}

func (s *GitInitialiserTestSuite) SetupTest() {
	t := s.T()
	s.pathWithGit = path.Join(getCurrentWorkingDirectory(), "/data/test-initialiser/git/exists")
	if !directoryExists(path.Join(s.pathWithGit, "/.git")) {
		if err := os.MkdirAll(path.Join(s.pathWithGit, "/.git"), os.ModePerm); err != nil {
			panic(err)
		}
	}
	s.pathWithoutGit = path.Join(getCurrentWorkingDirectory(), "/data/test-initialiser/git/non-existent")
	if directoryExists(path.Join(s.pathWithoutGit, "/.git")) {
		if err := removeDir(t, path.Join(s.pathWithoutGit, "/.git")); err != nil {
			panic(err)
		}
	}
	s.gitInitialiser = InitGitInitialiser(&GitInitialiserConfig{})
	s.gitInitialiser.logger.SetOutput(&s.logs)
}

func (s *GitInitialiserTestSuite) TearDownTest() {
	t := s.T()
	if directoryExists(path.Join(s.pathWithGit, "/.git")) {
		if err := removeDir(t, path.Join(s.pathWithGit, "/.git")); err != nil {
			panic(err)
		}
	}
	if directoryExists(path.Join(s.pathWithoutGit, "/.git")) {
		if err := removeDir(t, path.Join(s.pathWithoutGit, "/.git")); err != nil {
			panic(err)
		}
	}
}

func (s *GitInitialiserTestSuite) TestCheck() {
	t := s.T()
	s.gitInitialiser.Path = s.pathWithGit
	assert.True(t, s.gitInitialiser.Check())
	s.gitInitialiser.Path = s.pathWithoutGit
	assert.False(t, s.gitInitialiser.Check())
}

func (s *GitInitialiserTestSuite) TestConfirm() {
	reader := bufio.NewReader(strings.NewReader("y\n"))
	assert.True(s.T(), s.gitInitialiser.Confirm(reader))
	reader = bufio.NewReader(strings.NewReader("n\n"))
	assert.False(s.T(), s.gitInitialiser.Confirm(reader))
}

func (s *GitInitialiserTestSuite) TestHandle_skip() {
	s.gitInitialiser.Path = s.pathWithGit
	err := s.gitInitialiser.Handle(true)
	assert.Nil(s.T(), err)
	assert.Contains(
		s.T(),
		s.logs.String(),
		fmt.Sprintf(
			"skipping git repository initialisation at '%s'",
			s.pathWithGit,
		),
	)
}

func (s *GitInitialiserTestSuite) TestHandle_initialiseGit() {
	t := s.T()
	s.gitInitialiser.Path = s.pathWithoutGit
	err := s.gitInitialiser.Handle()
	assert.Nil(t, err)
}
