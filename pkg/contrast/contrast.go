package contrast

import (
	"fmt"
	"math"
)

type RGBColor struct {
	r uint8
	g uint8
	b uint8
}

func NewRGBColor(r, g, b uint8) *RGBColor {
	return &RGBColor{r, g, b}
}

// https://www.w3.org/TR/2008/REC-WCAG20-20081211/#relativeluminancedef
func getRelativeLuminance(l uint8) (rl float64) {
	base := uint8(255)
	c := float64(l) / float64(base)
	if c <= 0.03928 {
		rl = float64(c) / 12.92
	} else {
		//fmt.Println("else", math.Pow((float64(l)+0.055)/1.055, 2.4))
		rl = math.Pow((float64(c)+0.055)/1.055, 2.4)
	}

	return
}

// https://www.w3.org/TR/2008/REC-WCAG20-20081211/#contrastratio
func GetRatio(text, background *RGBColor) (ratio float64) {
	trrl := getRelativeLuminance(text.r)
	tgrl := getRelativeLuminance(text.g)
	tbrl := getRelativeLuminance(text.b)

	brrl := getRelativeLuminance(background.r)
	bgrl := getRelativeLuminance(background.g)
	bbrl := getRelativeLuminance(background.b)

	tl := trrl*0.2126 + tgrl*0.7152 + tbrl*0.0722
	bl := brrl*0.2126 + bgrl*0.7152 + bbrl*0.0722

	var lighter, darker float64
	if tl >= bl {
		lighter = tl
		darker = bl
	} else {
		lighter = bl
		darker = tl
	}
	fmt.Println("l, d", lighter, darker)
	return (lighter + 0.05) / (darker + 0.05)
}

func GetLevel(ratio float64) (level string) {
	if ratio >= 7.0 {
		level = "AAA"
	} else if ratio >= 4.5 {
		level = "AA"
	} else if ratio >= 3.0 {
		level = "A"
	} else {
		level = "-"
	}
	return
}
