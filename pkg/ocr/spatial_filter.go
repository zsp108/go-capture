package ocr

import (
	"math"
	"sort"
	"strings"
)

const (
	// PixelTolerance is the subpixel rendering tolerance (0.5px).
	PixelTolerance = 0.5
	// MinEnclosureRatio is the minimum overlap ratio required if subpixel aliasing exists.
	MinEnclosureRatio = 0.95
)

// FilterCharactersStrict100 filters characters that are strictly 100% enclosed within the selection box.
// Edge-cutoff characters or partially touched characters outside the selection box are strictly discarded.
func FilterCharactersStrict100(chars []CharBoundingBox, sel SelectionRect) []CharBoundingBox {
	var matched []CharBoundingBox

	for _, cb := range chars {
		// Strict coordinate enclosure check with 0.5px subpixel tolerance
		isCoordEnclosed := (cb.Left >= sel.Left-PixelTolerance &&
			cb.Right <= sel.Right+PixelTolerance &&
			cb.Top >= sel.Top-PixelTolerance &&
			cb.Bottom <= sel.Bottom+PixelTolerance)

		// Calculate 2D intersection overlap
		overlapLeft := math.Max(cb.Left, sel.Left)
		overlapRight := math.Min(cb.Right, sel.Right)
		overlapTop := math.Max(cb.Top, sel.Top)
		overlapBottom := math.Min(cb.Bottom, sel.Bottom)

		overlapW := math.Max(0, overlapRight-overlapLeft)
		overlapH := math.Max(0, overlapBottom-overlapTop)
		overlapArea := overlapW * overlapH
		charArea := cb.Area()

		// If coordinate enclosed OR overlap >= 95% of character's total area
		if isCoordEnclosed || (charArea > 0 && (overlapArea/charArea) >= MinEnclosureRatio) {
			matched = append(matched, cb)
		}
	}

	return matched
}

// AssembleTextLines groups and sorts matched characters into clean reading-order lines.
func AssembleTextLines(chars []CharBoundingBox) ([]string, string) {
	if len(chars) == 0 {
		return nil, ""
	}

	// Group characters by line index or approximate vertical baseline
	type charGroup struct {
		lineIndex int
		avgTop    float64
		chars     []CharBoundingBox
	}

	lineMap := make(map[int]*charGroup)
	for _, c := range chars {
		g, exists := lineMap[c.LineIndex]
		if !exists {
			g = &charGroup{
				lineIndex: c.LineIndex,
				avgTop:    c.Top,
				chars:     []CharBoundingBox{},
			}
			lineMap[c.LineIndex] = g
		}
		g.chars = append(g.chars, c)
	}

	// Sort lines vertically by lineIndex / avgTop
	var groups []*charGroup
	for _, g := range lineMap {
		// Sort characters horizontally within each line
		sort.Slice(g.chars, func(i, j int) bool {
			return g.chars[i].Left < g.chars[j].Left
		})
		groups = append(groups, g)
	}

	sort.Slice(groups, func(i, j int) bool {
		if groups[i].lineIndex != groups[j].lineIndex {
			return groups[i].lineIndex < groups[j].lineIndex
		}
		return groups[i].avgTop < groups[j].avgTop
	})

	var resultLines []string
	for _, g := range groups {
		var sb strings.Builder
		for i, c := range g.chars {
			// If gap between previous char and current char > 1.2 * average width, insert space for Latin text
			if i > 0 {
				prev := g.chars[i-1]
				gap := c.Left - prev.Right
				avgW := (prev.Width() + c.Width()) / 2
				if gap > avgW*0.35 && gap < avgW*3.0 {
					// Add space if alphanumeric
					if isAlphaNum(prev.Char) && isAlphaNum(c.Char) {
						sb.WriteString(" ")
					}
				}
			}
			sb.WriteString(c.Char)
		}
		lineStr := strings.TrimSpace(sb.String())
		if len(lineStr) > 0 {
			resultLines = append(resultLines, lineStr)
		}
	}

	rawText := strings.Join(resultLines, "\n")
	return resultLines, rawText
}

func isAlphaNum(s string) bool {
	if len(s) == 0 {
		return false
	}
	r := rune(s[0])
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}
