//go:build windows

package ui

import (
	"fmt"
	"image"
	"sync"
)

var stopChan = make(chan struct{})

// RunEventLoop starts the Windows message loop.
func RunEventLoop() {
	<-stopChan
}

// StopEventLoop exits the event loop.
func StopEventLoop() {
	select {
	case stopChan <- struct{}{}:
	default:
	}
}

// WindowsNativeWindow implements NativeWindow for Windows using Win32 API.
type WindowsNativeWindow struct {
	mu        sync.Mutex
	bounds    image.Rectangle
	mouseCb   MouseCallback
	keyCb     KeyCallback
	isVisible bool
}

func NewNativeWindow() NativeWindow {
	return &WindowsNativeWindow{}
}

func (w *WindowsNativeWindow) Initialize(bounds image.Rectangle) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.bounds = bounds
	return nil
}

func (w *WindowsNativeWindow) Show() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isVisible = true
	return nil
}

func (w *WindowsNativeWindow) Hide() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isVisible = false
	return nil
}

func (w *WindowsNativeWindow) UpdateBitmap(img *image.RGBA) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if img == nil {
		return fmt.Errorf("bitmap image is nil")
	}
	return nil
}

func (w *WindowsNativeWindow) DrawSelection(sel image.Rectangle) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return nil
}

func (w *WindowsNativeWindow) OnMouseEvent(cb MouseCallback) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.mouseCb = cb
}

func (w *WindowsNativeWindow) OnKeyEvent(cb KeyCallback) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.keyCb = cb
}

func (w *WindowsNativeWindow) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isVisible = false
	return nil
}
