package loupe

import (
	"image"
	"image/color"
	"testing"
)

func TestSampleOddGrid(t *testing.T) {
	// Create a test 100x100 image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// Paint center pixel (50, 50) with blue #3B82F6 -> (59, 130, 246)
	targetColor := color.RGBA{R: 59, G: 130, B: 246, A: 255}
	img.SetRGBA(50, 50, targetColor)

	// Sample 13x13 grid centered at (50, 50)
	snap := SampleOddGrid(img, 50, 50)

	// Verify dimensions
	if snap.CenterX != 50 || snap.CenterY != 50 {
		t.Fatalf("Expected center (50, 50), got (%d, %d)", snap.CenterX, snap.CenterY)
	}

	// Verify center pixel is at row=6, col=6
	centerCell := snap.Grid[6][6]
	if centerCell.X != 50 || centerCell.Y != 50 {
		t.Errorf("Expected cell (6,6) to map to pixel (50,50), got (%d, %d)", centerCell.X, centerCell.Y)
	}

	if centerCell.Color != targetColor {
		t.Errorf("Expected center color %v, got %v", targetColor, centerCell.Color)
	}

	if snap.CenterHex != "#3B82F6" {
		t.Errorf("Expected #3B82F6, got %s", snap.CenterHex)
	}
}

func TestColorSpaceConversions(t *testing.T) {
	c := color.RGBA{R: 239, G: 68, B: 68, A: 255} // #EF4444

	hex := RGBToHEX(c)
	if hex != "#EF4444" {
		t.Errorf("Expected #EF4444, got %s", hex)
	}

	parsed, err := ParseHEX(hex)
	if err != nil || parsed.R != 239 || parsed.G != 68 || parsed.B != 68 {
		t.Errorf("ParseHEX failed, got %v, err: %v", parsed, err)
	}

	h, s, l := RGBToHSL(c)
	if h < 0 || h > 360 || s < 0 || s > 100 || l < 0 || l > 100 {
		t.Errorf("Invalid HSL values: %.0f, %.0f, %.0f", h, s, l)
	}
}
