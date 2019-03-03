package main

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ColorTestSuite struct {
	suite.Suite
}

func TestColorTestSuite(t *testing.T) {
	suite.Run(t, new(ColorTestSuite))
}

func (s *ColorTestSuite) TestColor_NoOverwrite() {
	testString := Color("green", Color("blue", "HI")+" WORLD")
	assert.Contains(s.T(), testString, fmt.Sprintf("%s0m%s%vm", ColorStub, ColorStub, Palette["green"]))
}

func (s *ColorTestSuite) TestColorer() {
	t := s.T()
	assert.Contains(t, Colorer.Default("a"), strconv.Itoa(Palette["default"]))
	assert.Contains(t, Colorer.Bold("a"), strconv.Itoa(Palette["bold"]))
	assert.Contains(t, Colorer.Dim("a"), strconv.Itoa(Palette["dim"]))
	assert.Contains(t, Colorer.Underline("a"), strconv.Itoa(Palette["underline"]))
	assert.Contains(t, Colorer.Italics("a"), strconv.Itoa(Palette["italics"]))
	assert.Contains(t, Colorer.Black("a"), strconv.Itoa(Palette["black"]))
	assert.Contains(t, Colorer.Gray("a"), strconv.Itoa(Palette["gray"]))
	assert.Contains(t, Colorer.Grey("a"), strconv.Itoa(Palette["grey"]))
	assert.Contains(t, Colorer.Red("a"), strconv.Itoa(Palette["red"]))
	assert.Contains(t, Colorer.Yellow("a"), strconv.Itoa(Palette["yellow"]))
	assert.Contains(t, Colorer.Green("a"), strconv.Itoa(Palette["green"]))
	assert.Contains(t, Colorer.Cyan("a"), strconv.Itoa(Palette["cyan"]))
	assert.Contains(t, Colorer.Blue("a"), strconv.Itoa(Palette["blue"]))
	assert.Contains(t, Colorer.Violet("a"), strconv.Itoa(Palette["violet"]))
	assert.Contains(t, Colorer.LightGray("a"), strconv.Itoa(Palette["lgray"]))
	assert.Contains(t, Colorer.LightGrey("a"), strconv.Itoa(Palette["lgrey"]))
	assert.Contains(t, Colorer.LightRed("a"), strconv.Itoa(Palette["lred"]))
	assert.Contains(t, Colorer.LightYellow("a"), strconv.Itoa(Palette["lyellow"]))
	assert.Contains(t, Colorer.LightGreen("a"), strconv.Itoa(Palette["lgreen"]))
	assert.Contains(t, Colorer.LightCyan("a"), strconv.Itoa(Palette["lcyan"]))
	assert.Contains(t, Colorer.LightBlue("a"), strconv.Itoa(Palette["lblue"]))
	assert.Contains(t, Colorer.LightViolet("a"), strconv.Itoa(Palette["lviolet"]))
}
