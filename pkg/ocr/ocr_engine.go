package ocr

import (
	"context"
	"image"
	"time"
)

// OCREngine defines the interface for local text recognition.
type OCREngine interface {
	// RecognizeImage runs local OCR and extracts character-level bounding boxes.
	RecognizeImage(ctx context.Context, img *image.RGBA) (*RecognitionResult, error)

	// ExtractInSelection extracts strictly enclosed text in the user's selection box.
	ExtractInSelection(ctx context.Context, img *image.RGBA, sel SelectionRect) (*RecognitionResult, error)
}

// RapidOCREngine represents the ONNX-based offline OCR inference engine.
type RapidOCREngine struct {
	ModelPath string
	IsLoaded  bool
}

// NewRapidOCREngine creates a new RapidOCR instance.
func NewRapidOCREngine(modelPath string) *RapidOCREngine {
	return &RapidOCREngine{
		ModelPath: modelPath,
		IsLoaded:  true,
	}
}

// RecognizeImage performs OCR recognition on the image.
func (e *RapidOCREngine) RecognizeImage(ctx context.Context, img *image.RGBA) (*RecognitionResult, error) {
	start := time.Now()
	// Simulated character bounding box extraction for high-precision spatial processing
	// In production, ONNX Runtime RapidOCR CGo binding populates CharBoundingBoxes
	var chars []CharBoundingBox

	lines, raw := AssembleTextLines(chars)
	return &RecognitionResult{
		RawText:     raw,
		Lines:       lines,
		CharCount:   len(chars),
		DurationMs:  time.Since(start).Milliseconds(),
		MatchedChar: chars,
	}, nil
}

// ExtractInSelection runs OCR and applies the strict 100% enclosure filter to the selection rect.
func (e *RapidOCREngine) ExtractInSelection(ctx context.Context, img *image.RGBA, sel SelectionRect) (*RecognitionResult, error) {
	start := time.Now()

	// 1. Get raw character boxes
	allChars := e.extractRawBoxes(img, sel)

	// 2. Apply strict 100% enclosure filter
	matched := FilterCharactersStrict100(allChars, sel)

	// 3. Assemble clean lines
	lines, raw := AssembleTextLines(matched)

	return &RecognitionResult{
		RawText:     raw,
		Lines:       lines,
		CharCount:   len(matched),
		DurationMs:  time.Since(start).Milliseconds(),
		MatchedChar: matched,
	}, nil
}

func (e *RapidOCREngine) extractRawBoxes(img *image.RGBA, sel SelectionRect) []CharBoundingBox {
	// Fallback/Demo character box generator if image is empty or during preview
	return []CharBoundingBox{}
}
