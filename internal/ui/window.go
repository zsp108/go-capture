package ui

import (
	"image"
)

// MouseEventType defines native mouse events received from the OS window.
type MouseEventType int

const (
	MouseEventMouseDown MouseEventType = iota
	MouseEventMouseMove
	MouseEventMouseUp
	MouseEventDoubleClick
	MouseEventRightClick
)

// NativeMouseEvent represents a low-level OS mouse event.
type NativeMouseEvent struct {
	Type   MouseEventType
	X      int
	Y      int
	Button int
}

// KeyCallback defines the keyboard event callback handler.
type KeyCallback func(key string, ctrl bool, shift bool, meta bool)

// MouseCallback defines the mouse event callback handler.
type MouseCallback func(event NativeMouseEvent)

// NativeWindow defines the interface for pure OS-level frameless transparent overlay windows.
type NativeWindow interface {
	// Initialize creates a borderless topmost window matching display bounds.
	Initialize(bounds image.Rectangle) error

	// Show displays the fullscreen overlay.
	Show() error

	// Hide hides the overlay window.
	Hide() error

	// UpdateBitmap updates the displayed background bitmap.
	UpdateBitmap(img *image.RGBA) error

	// DrawSelection updates the cutout transparent region and border.
	DrawSelection(sel image.Rectangle) error

	// OnMouseEvent sets the mouse event listener callback.
	OnMouseEvent(cb MouseCallback)

	// OnKeyEvent sets the keyboard event listener callback.
	OnKeyEvent(cb KeyCallback)

	// Close destroys the native window handle.
	Close() error
}
