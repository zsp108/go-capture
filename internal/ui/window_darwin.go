//go:build darwin

package ui

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa -framework CoreGraphics -framework QuartzCore
#import <Cocoa/Cocoa.h>

@interface GoCaptureAppDelegate : NSObject <NSApplicationDelegate>
@end

@implementation GoCaptureAppDelegate
- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
    // Immediately responds to LaunchServices to prevent -1712 timeout error
    NSLog(@"[GoCapture] Application initialized and ready.");
}

- (NSApplicationTerminateReply)applicationShouldTerminate:(NSApplication *)sender {
    return NSTerminateNow;
}
@end

// Run the Cocoa event loop on the main OS thread
static void run_ns_app() {
    @autoreleasepool {
        NSApplication *app = [NSApplication sharedApplication];
        GoCaptureAppDelegate *delegate = [[GoCaptureAppDelegate alloc] init];
        [app setDelegate:delegate];
        [app setActivationPolicy:NSApplicationActivationPolicyAccessory];
        [app run];
    }
}

// Stop the Cocoa event loop safely
static void stop_ns_app() {
    @autoreleasepool {
        if (NSApp != nil) {
            [NSApp stop:nil];
            NSEvent* event = [NSEvent otherEventWithType:NSEventTypeApplicationDefined
                location:NSMakePoint(0, 0)
                modifierFlags:0
                timestamp:0
                windowNumber:0
                context:nil
                subtype:0
                data1:0
                data2:0];
            [NSApp postEvent:event atStart:YES];
        }
    }
}

// C-compatible bridge functions for Objective-C window management
static void* create_overlay_window(int x, int y, int w, int h) {
    @autoreleasepool {
        NSRect frame = NSMakeRect(x, y, w, h);
        NSWindow* window = [[NSWindow alloc] initWithContentRect:frame
            styleMask:NSWindowStyleMaskBorderless
            backing:NSBackingStoreBuffered
            defer:NO];

        [window setOpaque:NO];
        [window setBackgroundColor:[NSColor clearColor]];
        [window setLevel:NSPopUpMenuWindowLevel];
        [window setIgnoresMouseEvents:NO];
        [window setAcceptsMouseMovedEvents:YES];
        [window setCollectionBehavior:(NSWindowCollectionBehaviorCanJoinAllSpaces | NSWindowCollectionBehaviorFullScreenAuxiliary)];

        return (__bridge_retained void*)window;
    }
}

static void show_overlay_window(void* winPtr) {
    @autoreleasepool {
        if (winPtr != NULL) {
            NSWindow* window = (__bridge NSWindow*)winPtr;
            [window makeKeyAndOrderFront:nil];
            [NSApp activateIgnoringOtherApps:YES];
        }
    }
}

static void hide_overlay_window(void* winPtr) {
    @autoreleasepool {
        if (winPtr != NULL) {
            NSWindow* window = (__bridge NSWindow*)winPtr;
            [window orderOut:nil];
        }
    }
}

static void close_overlay_window(void* winPtr) {
    @autoreleasepool {
        if (winPtr != NULL) {
            NSWindow* window = (__bridge_transfer NSWindow*)winPtr;
            [window close];
        }
    }
}
*/
import "C"
import (
	"fmt"
	"image"
	"sync"
	"unsafe"
)

// RunEventLoop starts the native macOS Cocoa event loop, responding immediately to LaunchServices.
func RunEventLoop() {
	C.run_ns_app()
}

// StopEventLoop exits the Cocoa event loop.
func StopEventLoop() {
	C.stop_ns_app()
}

// DarwinNativeWindow implements NativeWindow for macOS using Cocoa NSWindow.
type DarwinNativeWindow struct {
	mu           sync.Mutex
	bounds       image.Rectangle
	windowHandle unsafe.Pointer
	mouseCb      MouseCallback
	keyCb        KeyCallback
	isVisible    bool
}

func NewNativeWindow() NativeWindow {
	return &DarwinNativeWindow{}
}

func (w *DarwinNativeWindow) Initialize(bounds image.Rectangle) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.bounds = bounds
	ptr := C.create_overlay_window(
		C.int(bounds.Min.X),
		C.int(bounds.Min.Y),
		C.int(bounds.Dx()),
		C.int(bounds.Dy()),
	)
	w.windowHandle = ptr
	return nil
}

func (w *DarwinNativeWindow) Show() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.windowHandle != nil {
		C.show_overlay_window(w.windowHandle)
	}
	w.isVisible = true
	return nil
}

func (w *DarwinNativeWindow) Hide() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.windowHandle != nil {
		C.hide_overlay_window(w.windowHandle)
	}
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

	if w.windowHandle != nil {
		C.close_overlay_window(w.windowHandle)
		w.windowHandle = nil
	}
	w.isVisible = false
	return nil
}
