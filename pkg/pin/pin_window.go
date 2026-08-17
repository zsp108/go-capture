package pin

import (
	"image"
	"math"
)

const (
	MinScale   = 0.20 // 20%
	MaxScale   = 3.00 // 300%
	ScaleStep  = 0.05 // 5% per wheel tick
	MinOpacity = 0.10 // 10%
	MaxOpacity = 1.00 // 100%
)

// PinnedWindow represents an always-on-top desktop pinned image window.
type PinnedWindow struct {
	ID         string      `json:"id"`
	X          float64     `json:"x"`
	Y          float64     `json:"y"`
	Width      float64     `json:"width"`
	Height     float64     `json:"height"`
	OriginalW  float64     `json:"original_w"`
	OriginalH  float64     `json:"original_h"`
	Scale      float64     `json:"scale"`    // 0.20 ~ 3.00
	Opacity    float64     `json:"opacity"`  // 0.10 ~ 1.00
	Rotation   int         `json:"rotation"` // 0, 90, 180, 270 degrees
	FlipX      bool        `json:"flip_x"`
	IsTopmost  bool        `json:"is_topmost"`
	Image      *image.RGBA `json:"-"`
}

// NewPinnedWindow creates a new pinned window instance.
func NewPinnedWindow(id string, x, y, w, h float64, img *image.RGBA) *PinnedWindow {
	return &PinnedWindow{
		ID:        id,
		X:         x,
		Y:         y,
		Width:     w,
		Height:    h,
		OriginalW: w,
		OriginalH: h,
		Scale:     1.0,
		Opacity:   1.0,
		Rotation:  0,
		FlipX:     false,
		IsTopmost: true,
		Image:     img,
	}
}

// AdjustScale updates the zoom scale with clamping (20% ~ 300%).
func (p *PinnedWindow) AdjustScale(delta float64) float64 {
	newScale := p.Scale + delta
	if newScale < MinScale {
		newScale = MinScale
	}
	if newScale > MaxScale {
		newScale = MaxScale
	}
	p.Scale = math.Round(newScale*100) / 100
	p.Width = p.OriginalW * p.Scale
	p.Height = p.OriginalH * p.Scale
	return p.Scale
}

// SetOpacity sets the window opacity with clamping (10% ~ 100%).
func (p *PinnedWindow) SetOpacity(opacity float64) {
	if opacity < MinOpacity {
		opacity = MinOpacity
	}
	if opacity > MaxOpacity {
		opacity = MaxOpacity
	}
	p.Opacity = opacity
}

// Rotate90 rotates the pinned image 90 degrees clockwise.
func (p *PinnedWindow) Rotate90() int {
	p.Rotation = (p.Rotation + 90) % 360
	return p.Rotation
}

// ToggleFlipX toggles the horizontal mirror state.
func (p *PinnedWindow) ToggleFlipX() bool {
	p.FlipX = !p.FlipX
	return p.FlipX
}
