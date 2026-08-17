package ocr

// SelectionRect defines the bounding coordinates of the user's selection box.
type SelectionRect struct {
	Left   float64 `json:"left"`
	Top    float64 `json:"top"`
	Right  float64 `json:"right"`
	Bottom float64 `json:"bottom"`
}

// Width returns the width of the selection rect.
func (r SelectionRect) Width() float64 {
	return r.Right - r.Left
}

// Height returns the height of the selection rect.
func (r SelectionRect) Height() float64 {
	return r.Bottom - r.Top
}

// CharBoundingBox represents a single Unicode character's spatial bounding box.
type CharBoundingBox struct {
	Char       string  `json:"char"`
	Left       float64 `json:"left"`
	Top        float64 `json:"top"`
	Right      float64 `json:"right"`
	Bottom     float64 `json:"bottom"`
	LineIndex  int     `json:"line_index"`
	Confidence float64 `json:"confidence"`
}

// Width returns the width of the character box.
func (cb CharBoundingBox) Width() float64 {
	return cb.Right - cb.Left
}

// Height returns the height of the character box.
func (cb CharBoundingBox) Height() float64 {
	return cb.Bottom - cb.Top
}

// Area returns the pixel area of the character box.
func (cb CharBoundingBox) Area() float64 {
	return cb.Width() * cb.Height()
}

// RecognitionResult represents the final assembled OCR output.
type RecognitionResult struct {
	RawText     string            `json:"raw_text"`
	Lines       []string          `json:"lines"`
	CharCount   int               `json:"char_count"`
	DurationMs  int64             `json:"duration_ms"`
	MatchedChar []CharBoundingBox `json:"matched_chars,omitempty"`
}
