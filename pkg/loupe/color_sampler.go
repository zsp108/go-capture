package loupe

import (
	"image"
	"image/color"
)

const (
	// GridSize is 13x13 odd matrix ensuring exact center alignment.
	GridSize = 13
	// HalfGrid is 6 pixels half-width.
	HalfGrid = 6
	// CenterIndex is the center row/col index (0-indexed: 6).
	CenterIndex = 6
)

// PixelSample represents a single sampled cell in the 13x13 matrix.
type PixelSample struct {
	X     int        `json:"x"`
	Y     int        `json:"y"`
	Row   int        `json:"row"`
	Col   int        `json:"col"`
	Color color.RGBA `json:"color"`
	Hex   string     `json:"hex"`
}

// LoupeSnapshot represents a 13x13 sampled snapshot around cursor (centerX, centerY).
type LoupeSnapshot struct {
	CenterX     int                  `json:"center_x"`
	CenterY     int                  `json:"center_y"`
	CenterColor color.RGBA           `json:"center_color"`
	CenterHex   string               `json:"center_hex"`
	CenterRGB   string               `json:"center_rgb"`
	CenterHSL   string               `json:"center_hsl"`
	Grid        [GridSize][GridSize]PixelSample `json:"grid"`
}

// SampleOddGrid samples a 13x13 odd matrix from the source image centered at (cx, cy).
func SampleOddGrid(src *image.RGBA, cx, cy int) *LoupeSnapshot {
	snap := &LoupeSnapshot{
		CenterX: cx,
		CenterY: cy,
	}

	bounds := src.Bounds()

	for row := 0; row < GridSize; row++ {
		for col := 0; col < GridSize; col++ {
			sampleX := cx + (col - HalfGrid)
			sampleY := cy + (row - HalfGrid)

			var c color.RGBA
			if sampleX >= bounds.Min.X && sampleX < bounds.Max.X &&
				sampleY >= bounds.Min.Y && sampleY < bounds.Max.Y {
				c = src.RGBAAt(sampleX, sampleY)
			} else {
				c = color.RGBA{R: 0, G: 0, B: 0, A: 255}
			}

			hexStr := RGBToHEX(c)
			sample := PixelSample{
				X:     sampleX,
				Y:     sampleY,
				Row:   row,
				Col:   col,
				Color: c,
				Hex:   hexStr,
			}
			snap.Grid[row][col] = sample

			if row == CenterIndex && col == CenterIndex {
				snap.CenterColor = c
				snap.CenterHex = hexStr
				snap.CenterRGB = RGBToString(c)
				h, s, l := RGBToHSL(c)
				snap.CenterHSL = HSLToString(h, s, l)
			}
		}
	}

	return snap
}
