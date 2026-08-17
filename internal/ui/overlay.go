package ui

import (
	"fmt"
	"image"
	"sync"
)

// OverlayController coordinates the native OS fullscreen overlay window and user interactions.
type OverlayController struct {
	mu           sync.Mutex
	nativeWindow NativeWindow
	currentImage *image.RGBA
	selection    image.Rectangle
	isActive     bool
}

// NewOverlayController creates an initialized native overlay controller.
func NewOverlayController() *OverlayController {
	return &OverlayController{
		nativeWindow: NewNativeWindow(),
	}
}

// StartSession initializes the native window with display dimensions and sets the frozen snapshot image.
func (c *OverlayController) StartSession(screenImg *image.RGBA) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if screenImg == nil {
		return fmt.Errorf("screen image is nil")
	}

	c.currentImage = screenImg
	bounds := screenImg.Bounds()

	if err := c.nativeWindow.Initialize(bounds); err != nil {
		return fmt.Errorf("failed to initialize native window: %w", err)
	}

	if err := c.nativeWindow.UpdateBitmap(screenImg); err != nil {
		return fmt.Errorf("failed to update bitmap: %w", err)
	}

	if err := c.nativeWindow.Show(); err != nil {
		return fmt.Errorf("failed to show native window: %w", err)
	}

	c.isActive = true
	return nil
}

// SetSelection updates the transparent cutout region on the native window.
func (c *OverlayController) SetSelection(rect image.Rectangle) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.selection = rect
	if c.isActive {
		return c.nativeWindow.DrawSelection(rect)
	}
	return nil
}

// Hide hides the overlay window.
func (c *OverlayController) Hide() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.isActive = false
	return c.nativeWindow.Hide()
}

// Close destroys the native overlay.
func (c *OverlayController) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.isActive = false
	return c.nativeWindow.Close()
}
