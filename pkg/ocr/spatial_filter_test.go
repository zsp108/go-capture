package ocr

import (
	"testing"
)

func TestFilterCharactersStrict100(t *testing.T) {
	// Sample simulated line: "Pin Windows to Top"
	// Bounding boxes:
	// 'P' [100, 50, 110, 70]
	// 'i' [110, 50, 116, 70]
	// 'n' [116, 50, 126, 70]
	// ' '
	// 'W' [134, 50, 148, 70]
	// 'i' [148, 50, 154, 70]
	// 'n' [154, 50, 164, 70]
	// 'd' [164, 50, 174, 70]
	// 'o' [174, 50, 184, 70]
	// 'w' [184, 50, 196, 70]
	// 's' [196, 50, 206, 70]
	// ' '
	// 't' [214, 50, 222, 70]
	// 'o' [222, 50, 232, 70]

	chars := []CharBoundingBox{
		{Char: "P", Left: 100, Top: 50, Right: 110, Bottom: 70, LineIndex: 0},
		{Char: "i", Left: 110, Top: 50, Right: 116, Bottom: 70, LineIndex: 0},
		{Char: "n", Left: 116, Top: 50, Right: 126, Bottom: 70, LineIndex: 0},
		{Char: "W", Left: 134, Top: 50, Right: 148, Bottom: 70, LineIndex: 0},
		{Char: "i", Left: 148, Top: 50, Right: 154, Bottom: 70, LineIndex: 0},
		{Char: "n", Left: 154, Top: 50, Right: 164, Bottom: 70, LineIndex: 0},
		{Char: "d", Left: 164, Top: 50, Right: 174, Bottom: 70, LineIndex: 0},
		{Char: "o", Left: 174, Top: 50, Right: 184, Bottom: 70, LineIndex: 0},
		{Char: "w", Left: 184, Top: 50, Right: 196, Bottom: 70, LineIndex: 0},
		{Char: "s", Left: 196, Top: 50, Right: 206, Bottom: 70, LineIndex: 0},
		{Char: "t", Left: 214, Top: 50, Right: 222, Bottom: 70, LineIndex: 0},
		{Char: "o", Left: 222, Top: 50, Right: 232, Bottom: 70, LineIndex: 0},
	}

	t.Run("Extract only 'Pin' word exactly", func(t *testing.T) {
		sel := SelectionRect{Left: 95, Top: 45, Right: 130, Bottom: 75}
		matched := FilterCharactersStrict100(chars, sel)
		lines, raw := AssembleTextLines(matched)
		if len(lines) != 1 || lines[0] != "Pin" {
			t.Errorf("Expected 'Pin', got: %v (raw: %s)", lines, raw)
		}
	})

	t.Run("Reject partially cutoff character on selection edge", func(t *testing.T) {
		// Selection right edge cuts 'd' in half at 168 (d is [164..174])
		// Characters 'P','i','n','W','i','n' should be matched, 'd' rejected!
		sel := SelectionRect{Left: 95, Top: 45, Right: 168, Bottom: 75}
		matched := FilterCharactersStrict100(chars, sel)
		_, raw := AssembleTextLines(matched)
		if raw != "Pin Win" {
			t.Errorf("Expected cutoff rejection resulting in 'Pin Win', got: %q", raw)
		}
	})

	t.Run("Extract 'Windows' only without 'Pin' or 'to'", func(t *testing.T) {
		sel := SelectionRect{Left: 130, Top: 45, Right: 210, Bottom: 75}
		matched := FilterCharactersStrict100(chars, sel)
		_, raw := AssembleTextLines(matched)
		if raw != "Windows" {
			t.Errorf("Expected 'Windows', got: %q", raw)
		}
	})
}
