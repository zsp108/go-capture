package app

import (
	"context"
	"fmt"
	"image"
	"sync"
	"time"

	"github.com/gocapture/go-capture/pkg/annotation"
	"github.com/gocapture/go-capture/pkg/capture"
	"github.com/gocapture/go-capture/pkg/clipboard"
	"github.com/gocapture/go-capture/pkg/hotkey"
	"github.com/gocapture/go-capture/pkg/loupe"
	"github.com/gocapture/go-capture/pkg/ocr"
	"github.com/gocapture/go-capture/pkg/pin"
)

// AppState defines the strict application state machine.
type AppState string

const (
	StateIdle       AppState = "IDLE"
	StateSelecting  AppState = "SELECTING"
	StateSelected   AppState = "SELECTED"
	StateResizing   AppState = "RESIZING"
	StateMoving     AppState = "MOVING"
	StateAnnotating AppState = "ANNOTATING"
	StatePinned     AppState = "PINNED"
)

// App is the core application coordinator.
type App struct {
	mu           sync.RWMutex
	state        AppState
	config       *Config
	capturer     capture.ScreenCapturer
	history      *annotation.HistoryManager
	pinManager   *pin.PinManager
	ocrEngine    ocr.OCREngine
	hotkeys      *hotkey.HotkeyListener

	// Current session snapshot
	currentScreen  *image.RGBA
	selectionRect  image.Rectangle
	activeTool     annotation.ShapeType
	isDropperMode  bool
	colorFormat    loupe.ColorFormat
}

// NewApp creates an initialized GoCapture application instance.
func NewApp(cfg *Config) *App {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	return &App{
		state:       StateIdle,
		config:      cfg,
		capturer:    capture.NewCapturer(),
		history:     annotation.NewHistoryManager(),
		pinManager:  pin.NewPinManager(),
		ocrEngine:   ocr.NewRapidOCREngine(""),
		hotkeys:     hotkey.NewHotkeyListener(),
		colorFormat: loupe.FormatHEX,
	}
}

// StartCapture transitions into SELECTING state and captures raw screen frame.
func (a *App) StartCapture() (*image.RGBA, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// 1. Capture screen at T0 snapshot
	img, err := a.capturer.CaptureDisplay(0)
	if err != nil {
		return nil, fmt.Errorf("screen capture failed: %w", err)
	}

	a.currentScreen = img
	a.state = StateSelecting
	a.history.Clear()
	a.activeTool = ""
	a.isDropperMode = false

	return img, nil
}

// SetSelection sets the active selection box and enters SELECTED state.
func (a *App) SetSelection(rect image.Rectangle) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.selectionRect = rect
	a.state = StateSelected
}

// ActivateAnnotationTool sets the active vector tool and transitions to ANNOTATING mode.
// In ANNOTATING mode, mouse events are exclusively consumed by drawing layer.
func (a *App) ActivateAnnotationTool(tool annotation.ShapeType) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.activeTool = tool
	a.state = StateAnnotating
}

// SamplePixelColor samples a 13x13 odd grid around (x, y).
func (a *App) SamplePixelColor(x, y int) *loupe.LoupeSnapshot {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.currentScreen == nil {
		return nil
	}
	return loupe.SampleOddGrid(a.currentScreen, x, y)
}

// ConfirmAndCopy synthesizes vector annotations and copies to system clipboard.
func (a *App) ConfirmAndCopy() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.currentScreen == nil || a.selectionRect.Empty() {
		return fmt.Errorf("no active selection")
	}

	// 1. Crop selection from base screen
	cropped := capture.CropImage(a.currentScreen, a.selectionRect)

	// 2. Render vector annotations
	shapes := a.history.GetShapes()
	rendered := annotation.RenderShapesOnImage(cropped, shapes)

	// 3. Push to system clipboard
	if err := clipboard.WritePNG(rendered); err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	// 4. Return to IDLE state
	a.state = StateIdle
	return nil
}

// PinToDesktop creates a pinned topmost window from current selection.
func (a *App) PinToDesktop() (*pin.PinnedWindow, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.currentScreen == nil || a.selectionRect.Empty() {
		return nil, fmt.Errorf("no active selection")
	}

	cropped := capture.CropImage(a.currentScreen, a.selectionRect)
	shapes := a.history.GetShapes()
	rendered := annotation.RenderShapesOnImage(cropped, shapes)

	w := float64(a.selectionRect.Dx())
	h := float64(a.selectionRect.Dy())
	x := float64(a.selectionRect.Min.X)
	y := float64(a.selectionRect.Min.Y)

	pinned := a.pinManager.CreatePin(x, y, w, h, rendered)
	a.state = StateIdle

	return pinned, nil
}

// ExtractOCR runs character-level 100% enclosure spatial extraction.
func (a *App) ExtractOCR(ctx context.Context) (*ocr.RecognitionResult, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.currentScreen == nil || a.selectionRect.Empty() {
		return nil, fmt.Errorf("no active selection")
	}

	sel := ocr.SelectionRect{
		Left:   float64(a.selectionRect.Min.X),
		Top:    float64(a.selectionRect.Min.Y),
		Right:  float64(a.selectionRect.Max.X),
		Bottom: float64(a.selectionRect.Max.Y),
	}

	res, err := a.ocrEngine.ExtractInSelection(ctx, a.currentScreen, sel)
	if err != nil {
		return nil, err
	}

	// Automatically copy to clipboard if configured
	if a.config.AutoCopy && len(res.RawText) > 0 {
		_ = clipboard.WriteText(res.RawText)
	}

	return res, nil
}

// GetState returns current application state.
func (a *App) GetState() AppState {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.state
}

// CancelCapture resets the session back to IDLE.
func (a *App) CancelCapture() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.state = StateIdle
	a.currentScreen = nil
	a.selectionRect = image.Rectangle{}
	a.history.Clear()
}

// History returns the annotation history manager.
func (a *App) History() *annotation.HistoryManager {
	return a.history
}

// PinManager returns the pinned window manager.
func (a *App) PinManager() *pin.PinManager {
	return a.pinManager
}

// InitHotkeys binds global shortcuts.
func (a *App) InitHotkeys(ctx context.Context) error {
	a.hotkeys.Register(hotkey.HotkeyCapture, func() {
		_, _ = a.StartCapture()
	})
	a.hotkeys.Register(hotkey.HotkeyCancel, func() {
		a.CancelCapture()
	})
	return a.hotkeys.Start(ctx)
}

// Shutdown cleanly stops workers.
func (a *App) Shutdown() {
	a.hotkeys.Stop()
	time.Sleep(10 * time.Millisecond)
}
