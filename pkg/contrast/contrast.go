package contrast

import (
	"gopkg.in/go-playground/colors.v1"
	"math"
)

// https://www.w3.org/TR/2008/REC-WCAG20-20081211/#relativeluminancedef
func getRelativeLuminance(l uint8) (rl float64) {
	base := uint8(255)
	c := float64(l) / float64(base)
	if c <= 0.03928 {
		rl = float64(c) / 12.92
	} else {
		rl = math.Pow((float64(c)+0.055)/1.055, 2.4)
	}

	return
}

// https://www.w3.org/TR/2008/REC-WCAG20-20081211/#contrastratio
func GetRatio(color1, color2 *colors.RGBColor) (ratio float64) {
	c1rrl := getRelativeLuminance(color1.R)
	c1grl := getRelativeLuminance(color1.G)
	c1brl := getRelativeLuminance(color1.B)

	c2rrl := getRelativeLuminance(color2.R)
	c2grl := getRelativeLuminance(color2.G)
	c2brl := getRelativeLuminance(color2.B)

	c1 := c1rrl*0.2126 + c1grl*0.7152 + c1brl*0.0722
	c2 := c2rrl*0.2126 + c2grl*0.7152 + c2brl*0.0722

	var lighter, darker float64
	if c1 >= c2 {
		lighter = c1
		darker = c2
	} else {
		lighter = c2
		darker = c1
	}

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
