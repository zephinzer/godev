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

func Color(color int, value string) string {
	return fmt.Sprintf("%s%vm%s%s0m", CFormat, color, value, CFormat)
}

func Multi(value string, colorCodes []int) string {
	retval := value + "\033[0m"
	for _, colorCode := range colorCodes {
		retval = Color(colorCode, retval)
	}

	return retval
}

func Default(value string) string {
	return Color(CDefault, value)
}

func Bold(value string) string {
	return Color(CBold, value)
}

func Dim(value string) string {
	return Color(CDim, value)
}

func Black(value string) string {
	return Color(CBlack, value)
}

func DarkGray(value string) string {
	return Color(CDarkGray, value)
}

func Red(value string) string {
	return Color(CRed, value)
}

func LightRed(value string) string {
	return Color(CLightRed, value)
}

func Green(value string) string {
	return Color(CGreen, value)
}

func LightGreen(value string) string {
	return Color(CLightGreen, value)
}

func Yellow(value string) string {
	return Color(CYellow, value)
}

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
