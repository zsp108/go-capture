package pin

import (
	"fmt"
	"image"
	"sync"
)

// PinManager coordinates all active desktop pinned windows.
type PinManager struct {
	mu      sync.RWMutex
	windows map[string]*PinnedWindow
	nextID  int
}

// NewPinManager creates an initialized PinManager.
func NewPinManager() *PinManager {
	return &PinManager{
		windows: make(map[string]*PinnedWindow),
		nextID:  1,
	}
}

// CreatePin creates and stores a new pinned window.
func (m *PinManager) CreatePin(x, y, w, h float64, img *image.RGBA) *PinnedWindow {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := fmt.Sprintf("pin_%d", m.nextID)
	m.nextID++

	win := NewPinnedWindow(id, x, y, w, h, img)
	m.windows[id] = win
	return win
}

// GetPin retrieves a pinned window by its ID.
func (m *PinManager) GetPin(id string) (*PinnedWindow, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	win, exists := m.windows[id]
	return win, exists
}

// ClosePin removes a pinned window.
func (m *PinManager) ClosePin(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.windows[id]; exists {
		delete(m.windows, id)
		return true
	}
	return false
}

// GetAllPins returns a list of all active pinned windows.
func (m *PinManager) GetAllPins() []*PinnedWindow {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*PinnedWindow, 0, len(m.windows))
	for _, win := range m.windows {
		result = append(result, win)
	}
	return result
}

// Count returns the number of active pinned windows.
func (m *PinManager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.windows)
}

// ClearAll removes all pinned windows.
func (m *PinManager) ClearAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.windows = make(map[string]*PinnedWindow)
}
