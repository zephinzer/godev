package main

import (
	"bufio"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UtilsTestSuite struct {
	suite.Suite
}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}

func (s *UtilsTestSuite) TestConfigCommaDelimitedStringSet() {
	ccds := ConfigCommaDelimitedString{"a", "b"}
	ccds.Set("c,d")
	assert.Len(s.T(), ccds, 4, "'c' and 'd' should have been separated items")
	ccds.Set("e f")
	assert.Len(s.T(), ccds, 5, "'e' and 'f' should not have been separated")
}

func (s *UtilsTestSuite) TestConfigCommaDelimitedStringString() {
	ccds := ConfigCommaDelimitedString{"a", "b", "c", "d"}
	assert.Len(s.T(), ccds.String(), 7)
	assert.Equal(s.T(), 3, strings.Count(ccds.String(), ","))
}

func (s *UtilsTestSuite) TestConfigMultiflagStringSet() {
	ccds := ConfigMultiflagString{"a", "b"}
	ccds.Set("c,d")
	assert.Len(s.T(), ccds, 3, "'c' and 'd' should not have been separated items")
}

func (s *UtilsTestSuite) TestConfigMultiflagStringString() {
	ccds := ConfigMultiflagString{"a", "b", "c", "d"}
	assert.Len(s.T(), ccds.String(), 7)
	assert.Equal(s.T(), 3, strings.Count(ccds.String(), ","))
}

func (s *UtilsTestSuite) Test_confirm_withReply() {
	assert.True(s.T(), confirm(bufio.NewReader(strings.NewReader("y\n")), "hi", true))
	assert.False(s.T(), confirm(bufio.NewReader(strings.NewReader("n\n")), "hi", true))
}

func (s *UtilsTestSuite) Test_confirm_withWindowsReply() {
	assert.True(s.T(), confirm(bufio.NewReader(strings.NewReader("y\r\n")), "hi", true))
	assert.False(s.T(), confirm(bufio.NewReader(strings.NewReader("n\r\n")), "hi", true))
}

func (s *UtilsTestSuite) Test_confirm_withWeirdReplyNoRetry() {
	assert.False(s.T(), confirm(bufio.NewReader(strings.NewReader("something\n")), "hi", true))
	assert.False(s.T(), confirm(bufio.NewReader(strings.NewReader("something\n")), "hi", true))
}

func (s *UtilsTestSuite) Test_confirm_withWeirdReplyAndRetry() {
	assert.True(s.T(), confirm(bufio.NewReader(strings.NewReader("something\ny\n")), "hi", true, "retry please"))
	assert.False(s.T(), confirm(bufio.NewReader(strings.NewReader("something\nn\n")), "hi", true, "retry please"))
}

func (s *UtilsTestSuite) Test_confirm_withoutReply() {
	assert.True(s.T(), confirm(bufio.NewReader(strings.NewReader("\n")), "hi", true))
	assert.False(s.T(), confirm(bufio.NewReader(strings.NewReader("\n")), "hi", false))
}

func (s *UtilsTestSuite) Test_directoryExists() {
	assert.True(s.T(), directoryExists(getCurrentWorkingDirectory()))
}

func (s *UtilsTestSuite) Test_directoryExists_onFile() {
	assert.False(s.T(), directoryExists(path.Join(getCurrentWorkingDirectory(), "/main.go")))
}

func (s *UtilsTestSuite) Test_fileExists() {
	assert.True(s.T(), fileExists(path.Join(getCurrentWorkingDirectory(), "/main.go")))
}

func (s *UtilsTestSuite) Test_fileExists_onDirectory() {
	assert.False(s.T(), fileExists(getCurrentWorkingDirectory()))
}
