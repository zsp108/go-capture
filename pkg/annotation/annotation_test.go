package annotation

import (
	"image"
	"image/color"
	"testing"
)

func TestHistoryManagerUndoRedo(t *testing.T) {
	hm := NewHistoryManager()

	if hm.CanUndo() || hm.CanRedo() {
		t.Errorf("Initial state should not allow undo/redo")
	}

	shape1 := &AnnotationShape{
		ID:        "1",
		Type:      ToolRect,
		Color:     color.RGBA{239, 68, 68, 255},
		StrokeWidth: 2,
		StartPoint: Point2D{X: 10, Y: 10},
		EndPoint:   Point2D{X: 50, Y: 50},
	}

	shape2 := &AnnotationShape{
		ID:        "2",
		Type:      ToolStep,
		Color:     color.RGBA{59, 130, 246, 255},
		StartPoint: Point2D{X: 30, Y: 30},
		StepIndex: 1,
	}

	hm.Push(shape1)
	hm.Push(shape2)

	if len(hm.GetShapes()) != 2 {
		t.Fatalf("Expected 2 shapes, got %d", len(hm.GetShapes()))
	}
	if hm.GetNextStepIndex() != 2 {
		t.Errorf("Expected next step index 2, got %d", hm.GetNextStepIndex())
	}

	// Undo shape2
	popped := hm.Undo()
	if popped == nil || popped.ID != "2" {
		t.Errorf("Expected popped shape 2, got %v", popped)
	}
	if len(hm.GetShapes()) != 1 {
		t.Errorf("Expected 1 shape remaining after undo, got %d", len(hm.GetShapes()))
	}
	if !hm.CanRedo() {
		t.Errorf("Should be able to redo after undo")
	}

	// Redo shape2
	restored := hm.Redo()
	if restored == nil || restored.ID != "2" {
		t.Errorf("Expected restored shape 2, got %v", restored)
	}
	if len(hm.GetShapes()) != 2 {
		t.Errorf("Expected 2 shapes after redo, got %d", len(hm.GetShapes()))
	}
}

func TestRenderShapesOnImage(t *testing.T) {
	base := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			base.SetRGBA(x, y, color.RGBA{0, 0, 0, 255})
		}
	}

	rectShape := &AnnotationShape{
		ID:        "r1",
		Type:      ToolRect,
		Color:     color.RGBA{255, 0, 0, 255},
		StrokeWidth: 2,
		StartPoint: Point2D{X: 20, Y: 20},
		EndPoint:   Point2D{X: 80, Y: 80},
	}

	out := RenderShapesOnImage(base, []*AnnotationShape{rectShape})
	// Check border pixel is painted red
	c := out.RGBAAt(20, 20)
	if c.R != 255 || c.G != 0 || c.B != 0 {
		t.Errorf("Expected pixel (20,20) to be red, got: %v", c)
	}
}
