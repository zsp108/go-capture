//go:build linux

package ui

import (
	"fmt"
	"image"
	"sync"
)

var stopChan = make(chan struct{})

// RunEventLoop starts the Linux event loop.
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

// LinuxNativeWindow implements NativeWindow for Linux using X11 / Wayland.
type LinuxNativeWindow struct {
	mu        sync.Mutex
	bounds    image.Rectangle
	mouseCb   MouseCallback
	keyCb     KeyCallback
	isVisible bool
}

func NewNativeWindow() NativeWindow {
	return &LinuxNativeWindow{}
}

func (w *LinuxNativeWindow) Initialize(bounds image.Rectangle) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.bounds = bounds
	return nil
}

func (w *LinuxNativeWindow) Show() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isVisible = true
	return nil
}

func (w *LinuxNativeWindow) Hide() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isVisible = false
	return nil
}

func (w *LinuxNativeWindow) UpdateBitmap(img *image.RGBA) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if img == nil {
		return fmt.Errorf("bitmap image is nil")
	}
	return nil
}

func (w *LinuxNativeWindow) DrawSelection(sel image.Rectangle) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return nil
}

func (w *LinuxNativeWindow) OnMouseEvent(cb MouseCallback) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.mouseCb = cb
}

func (w *LinuxNativeWindow) OnKeyEvent(cb KeyCallback) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.keyCb = cb
}

func (w *LinuxNativeWindow) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isVisible = false
	return nil
}
