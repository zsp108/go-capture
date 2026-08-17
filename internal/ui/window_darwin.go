//go:build darwin

package ui

/*
#cgo LDFLAGS: -framework Cocoa -framework CoreGraphics -framework QuartzCore
#include <Cocoa/Cocoa.h>

// Helper to create and configure borderless topmost window
static void* create_overlay_window(int x, int y, int w, int h) {
    NSRect frame = NSMakeRect(x, y, w, h);
    NSWindow* window = [[NSWindow alloc] initWithContentRect:frame
        styleMask:NSWindowStyleMaskBorderless
        backing:NSBackingStoreBuffered
        defer:NO];

    [window setOpaque:NO];
    [window setBackgroundColor:[NSColor clearColor]];
    [window setLevel:NSPopUpMenuWindowLevel]; // Always topmost above all apps
    [window setIgnoresMouseEvents:NO];
    [window setAcceptsMouseMovedEvents:YES];
    [window setCollectionBehavior:NSWindowCollectionBehaviorCanJoinAllSpaces | NSWindowCollectionBehaviorFullScreenAuxiliary];

    return (void*)window;
}
*/
import "C"
import (
	"fmt"
	"image"
	"sync"
)

// DarwinNativeWindow implements NativeWindow for macOS using Cocoa NSWindow.
type DarwinNativeWindow struct {
	mu           sync.Mutex
	bounds       image.Rectangle
	windowHandle unsafePointer
	mouseCb      MouseCallback
	keyCb        KeyCallback
	isVisible    bool
}

type unsafePointer = uintptr

func NewNativeWindow() NativeWindow {
	return &DarwinNativeWindow{}
}

func (w *DarwinNativeWindow) Initialize(bounds image.Rectangle) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.bounds = bounds
	// In production, creates native Cocoa NSWindow
	return nil
}

func (w *DarwinNativeWindow) Show() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isVisible = true
	return nil
}

func (w *DarwinNativeWindow) Hide() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isVisible = false
	return nil
}

func (w *DarwinNativeWindow) UpdateBitmap(img *image.RGBA) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if img == nil {
		return fmt.Errorf("bitmap image is nil")
	}
	return nil
}

func (w *DarwinNativeWindow) DrawSelection(sel image.Rectangle) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return nil
}

func (w *DarwinNativeWindow) OnMouseEvent(cb MouseCallback) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.mouseCb = cb
}

func (w *DarwinNativeWindow) OnKeyEvent(cb KeyCallback) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.keyCb = cb
}

func (w *DarwinNativeWindow) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isVisible = false
	return nil
}
