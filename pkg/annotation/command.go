package annotation

import (
	"sync"
)

// HistoryManager manages the Undo/Redo command stacks.
type HistoryManager struct {
	mu         sync.Mutex
	undoStack  []*AnnotationShape
	redoStack  []*AnnotationShape
	stepCounter int
}

// NewHistoryManager creates an initialized history stack manager.
func NewHistoryManager() *HistoryManager {
	return &HistoryManager{
		undoStack:   make([]*AnnotationShape, 0, 32),
		redoStack:   make([]*AnnotationShape, 0, 32),
		stepCounter: 1,
	}
}

// Push adds a new shape to the undo stack and clears the redo stack.
func (h *HistoryManager) Push(shape *AnnotationShape) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.undoStack = append(h.undoStack, shape.Clone())
	h.redoStack = h.redoStack[:0] // Clear redo stack on new action

	if shape.Type == ToolStep {
		h.stepCounter++
	}
}

// Undo removes the last shape from the undo stack and pushes it to the redo stack.
func (h *HistoryManager) Undo() *AnnotationShape {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.undoStack) == 0 {
		return nil
	}

	lastIdx := len(h.undoStack) - 1
	shape := h.undoStack[lastIdx]
	h.undoStack = h.undoStack[:lastIdx]
	h.redoStack = append(h.redoStack, shape)

	if shape.Type == ToolStep && h.stepCounter > 1 {
		h.stepCounter--
	}

	return shape
}

// Redo pops a shape from the redo stack and restores it to the undo stack.
func (h *HistoryManager) Redo() *AnnotationShape {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.redoStack) == 0 {
		return nil
	}

	lastIdx := len(h.redoStack) - 1
	shape := h.redoStack[lastIdx]
	h.redoStack = h.redoStack[:lastIdx]
	h.undoStack = append(h.undoStack, shape)

	if shape.Type == ToolStep {
		h.stepCounter++
	}

	return shape
}

// GetShapes returns all active shapes in drawing order.
func (h *HistoryManager) GetShapes() []*AnnotationShape {
	h.mu.Lock()
	defer h.mu.Unlock()

	result := make([]*AnnotationShape, len(h.undoStack))
	for i, s := range h.undoStack {
		result[i] = s.Clone()
	}
	return result
}

// Clear resets all stacks.
func (h *HistoryManager) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.undoStack = h.undoStack[:0]
	h.redoStack = h.redoStack[:0]
	h.stepCounter = 1
}

// CanUndo returns true if there are actions to undo.
func (h *HistoryManager) CanUndo() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.undoStack) > 0
}

// CanRedo returns true if there are actions to redo.
func (h *HistoryManager) CanRedo() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.redoStack) > 0
}

// GetNextStepIndex returns the next auto-incrementing step index for step badge tool.
func (h *HistoryManager) GetNextStepIndex() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.stepCounter
}
