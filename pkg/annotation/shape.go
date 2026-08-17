package annotation

import (
	"image/color"
)

// ShapeType defines the vector annotation tool types.
type ShapeType string

const (
	ToolRect        ShapeType = "rect"
	ToolArrow       ShapeType = "arrow"
	ToolPen         ShapeType = "pen"
	ToolHighlighter ShapeType = "highlighter"
	ToolMosaic      ShapeType = "mosaic"
	ToolStep        ShapeType = "step"
	ToolText        ShapeType = "text"
)

// Point2D represents a 2D coordinate in selection space.
type Point2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// AnnotationShape represents a single vector drawing entity.
type AnnotationShape struct {
	ID        string      `json:"id"`
	Type      ShapeType   `json:"type"`
	Color     color.RGBA  `json:"color"`
	ColorHex  string      `json:"color_hex"`
	StrokeWidth float64   `json:"stroke_width"`
	Opacity   float64     `json:"opacity"` // 0.1 ~ 1.0

	// Geometry points
	StartPoint Point2D   `json:"start_point"`
	EndPoint   Point2D   `json:"end_point"`
	Points     []Point2D `json:"points,omitempty"` // For Pen / Highlighter path

	// Properties for Step badge & Text
	StepIndex int    `json:"step_index,omitempty"`
	Text      string `json:"text,omitempty"`
	FontSize  int    `json:"font_size,omitempty"`

	// Mosaic block size
	BlockSize int `json:"block_size,omitempty"`
}

// Clone creates a deep copy of the annotation shape.
func (s *AnnotationShape) Clone() *AnnotationShape {
	cp := *s
	if len(s.Points) > 0 {
		cp.Points = make([]Point2D, len(s.Points))
		copy(cp.Points, s.Points)
	}
	return &cp
}
