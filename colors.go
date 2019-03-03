package main

import (
	"fmt"
	"strconv"
	"strings"
)

// ColorStub is the text string stub which the shell
// can interpret as an ANSI color
const ColorStub = "\033["

// Colorer is the convenience function to color things
var Colorer = Colors{}

// Palette stores the color code to int map
var Palette = map[string]int{
	"default":   0,
	"bold":      1,
	"dim":       2,
	"italics":   3,
	"underline": 4,
	"black":     30,
	"white":     97,
	"gray":      90,
	"grey":      90,
	"red":       31,
	"yellow":    33,
	"green":     32,
	"cyan":      36,
	"blue":      34,
	"violet":    35,
	"bgwhite":   47,
	"bgblack":   40,
	"bgred":     41,
	"bgyellow":  43,
	"bggreen":   42,
	"bgcyan":    46,
	"bgblue":    44,
	"bgviolet":  45,
	"lgray":     37,
	"lgrey":     37,
	"lred":      91,
	"lyellow":   93,
	"lgreen":    92,
	"lcyan":     96,
	"lblue":     94,
	"lviolet":   95,
}

// Color applies the color code :color to the string :value
func Color(color string, value string) string {
	format := ColorStub + strconv.Itoa(Palette[color]) + "m"
	unformat := ColorStub + "0m"
	finalValue := value
	if strings.Contains(value, unformat) {
		finalValue = strings.Replace(value, unformat, unformat+format, -1)
	}
	return fmt.Sprintf("%s%s%s%s", ColorStub, format, finalValue, unformat)
}

// Colors defines a struct for the Colorer
type Colors struct{}

// Default sets the string :value to the default color
func (*Colors) Default(value string) string {
	return Color("default", value)
}

// Bold sets the string :value to bold
func (*Colors) Bold(value string) string {
	return Color("bold", value)
}

// Dim dims the string :value
func (*Colors) Dim(value string) string {
	return Color("dim", value)
}

// Italics italicises the string :value
func (*Colors) Italics(value string) string {
	return Color("italics", value)
}

// Underline underlines the string :value
func (*Colors) Underline(value string) string {
	return Color("underline", value)
}

// Black sets the string :value to black
func (*Colors) Black(value string) string {
	return Color("black", value)
}

// Gray sets the string :value to gray
func (*Colors) Gray(value string) string {
	return Color("gray", value)
}

// Grey sets the string :value to grey
func (*Colors) Grey(value string) string {
	return Color("grey", value)
}

// Red sets the string :value to red
func (*Colors) Red(value string) string {
	return Color("red", value)
}

// LightRed sets the string :value to light red
func (*Colors) LightRed(value string) string {
	return Color("lred", value)
}

// Green sets the string :value to green
func (*Colors) Green(value string) string {
	return Color("green", value)
}

// LightGreen sets the string :value to light green
func (*Colors) LightGreen(value string) string {
	return Color("lgreen", value)
}

// Yellow sets the string :value to yellow
func (*Colors) Yellow(value string) string {
	return Color("yellow", value)
}

// LightYellow sets the string :value to light yellow
func (*Colors) LightYellow(value string) string {
	return Color("lyellow", value)
}

// Blue sets the string :value to blue
func (*Colors) Blue(value string) string {
	return Color("blue", value)
}

// LightBlue sets the string :value to light blue
func (*Colors) LightBlue(value string) string {
	return Color("lblue", value)
}

// Violet sets the string :value to violet
func (*Colors) Violet(value string) string {
	return Color("violet", value)
}

// LightViolet sets the string :value to light violet
func (*Colors) LightViolet(value string) string {
	return Color("lviolet", value)
}

// Cyan sets the string :value to cyan
func (*Colors) Cyan(value string) string {
	return Color("cyan", value)
}

// LightCyan sets the string :value to light cyan
func (*Colors) LightCyan(value string) string {
	return Color("lcyan", value)
}

// LightGray sets the string :value to light gray
func (*Colors) LightGray(value string) string {
	return Color("lgray", value)
}

// LightGrey sets the string :value to light grey
func (*Colors) LightGrey(value string) string {
	return Color("lgrey", value)
}

// White sets the string :value to light white
func (*Colors) White(value string) string {
	return Color("white", value)
}
