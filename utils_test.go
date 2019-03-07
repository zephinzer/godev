package main

import (
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
