package main

import "fmt"

const (
	// CDefault - color code
	CDefault = 0
	// CBold - color code
	CBold = 1
	// CDim - color code
	CDim = 2
	// CBlack - color code
	CBlack = 30
	// CDarkGray - color code
	CDarkGray = 90
	// CRed - color code
	CRed = 31
	// CLightRed - color code
	CLightRed = 91
	// CGreen - color code
	CGreen = 32
	// CLightGreen - color code
	CLightGreen = 92
	// CYellow - color code
	CYellow = 33
	// CLightYellow - color code
	CLightYellow = 93
	// CBlue - color code
	CBlue = 34
	// CLightBlue - color code
	CLightBlue = 94
	// CViolet - color code
	CViolet = 35
	// CLightViolet - color code
	CLightViolet = 95
	// CCyan - color code
	CCyan = 36
	// CLightCyan - color code
	CLightCyan = 96
	// CLightGray - color code
	CLightGray = 37
	// CWhite - color code
	CWhite = 97
	// CFormat - color code
	CFormat = "\033["
)

// Colors applies the color code :color to the string :value
func Color(color int, value string) string {
	return fmt.Sprintf("%s%vm%s%s0m", CFormat, color, value, CFormat)
}

// Colors applies multiple :colorCodes to the string :value
func Colors(colorCodes []int, value string) string {
	retval := value + "\033[0m"
	for _, colorCode := range colorCodes {
		retval = Color(colorCode, retval)
	}

	return retval
}

// Default sets the string :value to the default color
func Default(value string) string {
	return Color(CDefault, value)
}

// Bold sets the string :value to bold
func Bold(value string) string {
	return Color(CBold, value)
}

// Dim dims the string :value
func Dim(value string) string {
	return Color(CDim, value)
}

// Black sets the string :value to black
func Black(value string) string {
	return Color(CBlack, value)
}

// DarkGray sets the string :value to dark gray
func DarkGray(value string) string {
	return Color(CDarkGray, value)
}

// Red sets the string :value to red
func Red(value string) string {
	return Color(CRed, value)
}

// LightRed sets the string :value to light red
func LightRed(value string) string {
	return Color(CLightRed, value)
}

// Green sets the string :value to green
func Green(value string) string {
	return Color(CGreen, value)
}

// LightGreen sets the string :value to light green
func LightGreen(value string) string {
	return Color(CLightGreen, value)
}

// Yellow sets the string :value to yellow
func Yellow(value string) string {
	return Color(CYellow, value)
}

// LightYellow sets the string :value to light yellow
func LightYellow(value string) string {
	return Color(CLightYellow, value)
}

func Blue(value string) string {
	return Color(CBlue, value)
}

func LightBlue(value string) string {
	return Color(CLightBlue, value)
}

func Violet(value string) string {
	return Color(CViolet, value)
}

func LightViolet(value string) string {
	return Color(CLightViolet, value)
}

func Cyan(value string) string {
	return Color(CCyan, value)
}

func LightCyan(value string) string {
	return Color(CLightCyan, value)
}

func LightGray(value string) string {
	return Color(CLightGray, value)
}

func White(value string) string {
	return Color(CWhite, value)
}
