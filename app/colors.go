package main

import "fmt"

const (
	CDefault     = 0
	CBold        = 1
	CDim         = 2
	CBlack       = 30
	CDarkGray    = 90
	CRed         = 31
	CLightRed    = 91
	CGreen       = 32
	CLightGreen  = 92
	CYellow      = 33
	CLightYellow = 93
	CBlue        = 34
	CLightBlue   = 94
	CViolet      = 35
	CLightViolet = 95
	CCyan        = 36
	CLightCyan   = 96
	CLightGray   = 37
	CWhite       = 97
	CFormat      = "\033["
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
