package loupe

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
)

// ColorFormat represents the supported color string formats.
type ColorFormat string

const (
	FormatHEX ColorFormat = "HEX"
	FormatRGB ColorFormat = "RGB"
	FormatHSL ColorFormat = "HSL"
)

// RGBToHEX converts an RGBA color to a standard 6-digit hex string (#RRGGBB).
func RGBToHEX(c color.RGBA) string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

// RGBToString formats an RGBA color to "rgb(r, g, b)".
func RGBToString(c color.RGBA) string {
	return fmt.Sprintf("rgb(%d, %d, %d)", c.R, c.G, c.B)
}

// RGBToHSL converts an RGBA color to HSL (Hue 0..360, Saturation 0..100%, Lightness 0..100%).
func RGBToHSL(c color.RGBA) (h float64, s float64, l float64) {
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	delta := max - min

	l = (max + min) / 2.0

	if delta == 0 {
		h = 0
		s = 0
	} else {
		if l < 0.5 {
			s = delta / (max + min)
		} else {
			s = delta / (2.0 - max - min)
		}

		if max == r {
			h = (g - b) / delta
			if g < b {
				h += 6.0
			}
		} else if max == g {
			h = (b-r)/delta + 2.0
		} else {
			h = (r-g)/delta + 4.0
		}
		h *= 60.0
	}

	return math.Round(h), math.Round(s * 100.0), math.Round(l * 100.0)
}

// HSLToString formats HSL values to "hsl(h, s%, l%)".
func HSLToString(h, s, l float64) string {
	return fmt.Sprintf("hsl(%.0f, %.0f%%, %.0f%%)", h, s, l)
}

// FormatColor returns the color formatted in the specified format.
func FormatColor(c color.RGBA, format ColorFormat) string {
	switch format {
	case FormatRGB:
		return RGBToString(c)
	case FormatHSL:
		h, s, l := RGBToHSL(c)
		return HSLToString(h, s, l)
	case FormatHEX:
		fallthrough
	default:
		return RGBToHEX(c)
	}
}

// ParseHEX parses a hex string (e.g., "#3B82F6" or "3B82F6") to color.RGBA.
func ParseHEX(hexStr string) (color.RGBA, error) {
	hexStr = strings.TrimPrefix(hexStr, "#")
	if len(hexStr) == 3 {
		hexStr = string([]byte{hexStr[0], hexStr[0], hexStr[1], hexStr[1], hexStr[2], hexStr[2]})
	}
	if len(hexStr) != 6 {
		return color.RGBA{0, 0, 0, 255}, fmt.Errorf("invalid hex string length: %s", hexStr)
	}

	val, err := strconv.ParseUint(hexStr, 16, 32)
	if err != nil {
		return color.RGBA{0, 0, 0, 255}, err
	}

	return color.RGBA{
		R: uint8(val >> 16),
		G: uint8((val >> 8) & 0xFF),
		B: uint8(val & 0xFF),
		A: 255,
	}, nil
}
