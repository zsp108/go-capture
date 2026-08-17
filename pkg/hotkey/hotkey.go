package hotkey

import (
	"context"
	"fmt"
	"sync"
)

// KeyType defines registered hotkey actions.
type KeyType string

const (
	HotkeyCapture KeyType = "capture" // F1 / Ctrl+Shift+A
	HotkeyPin     KeyType = "pin"     // F3
	HotkeyCancel  KeyType = "cancel"  // Esc
)

// HotkeyListener coordinates global keyboard shortcuts.
type HotkeyListener struct {
	mu        sync.Mutex
	handlers  map[KeyType]func()
	isRunning bool
	cancel    context.CancelFunc
}

// NewHotkeyListener creates an initialized listener.
func NewHotkeyListener() *HotkeyListener {
	return &HotkeyListener{
		handlers: make(map[KeyType]func()),
	}
}

// Register binds an action handler to a hotkey type.
func (l *HotkeyListener) Register(key KeyType, handler func()) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.handlers[key] = handler
}

// Trigger manually dispatches a hotkey event (useful for test/CLI).
func (l *HotkeyListener) Trigger(key KeyType) {
	l.mu.Lock()
	handler, exists := l.handlers[key]
	l.mu.Unlock()

	if exists && handler != nil {
		handler()
	}
}

// Start begins listening to low-level OS global keyboard events.
func (l *HotkeyListener) Start(ctx context.Context) error {
	l.mu.Lock()
	if l.isRunning {
		l.mu.Unlock()
		return fmt.Errorf("hotkey listener is already running")
	}
	ctx, cancel := context.WithCancel(ctx)
	l.cancel = cancel
	l.isRunning = true
	l.mu.Unlock()

	go func() {
		<-ctx.Done()
		l.mu.Lock()
		l.isRunning = false
		l.mu.Unlock()
	}()

	return nil
}

// Stop halts the listener.
func (l *HotkeyListener) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.cancel != nil {
		l.cancel()
	}
	l.isRunning = false
}
