package pin

import (
	"image"
	"testing"
)

func TestPinnedWindowTransformations(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 400, 300))
	win := NewPinnedWindow("pin_1", 100, 100, 400, 300, img)

	// Test scale up
	win.AdjustScale(0.10)
	if win.Scale != 1.10 || win.Width != 440 || win.Height != 330 {
		t.Errorf("Scale adjustment failed, got scale: %f, w: %f, h: %f", win.Scale, win.Width, win.Height)
	}

	// Test max scale clamping
	win.AdjustScale(5.0)
	if win.Scale != MaxScale {
		t.Errorf("Expected max scale %f, got %f", MaxScale, win.Scale)
	}

	// Test min scale clamping
	win.AdjustScale(-5.0)
	if win.Scale != MinScale {
		t.Errorf("Expected min scale %f, got %f", MinScale, win.Scale)
	}

	// Test rotation
	r1 := win.Rotate90()
	if r1 != 90 {
		t.Errorf("Expected 90 degrees, got %d", r1)
	}
	r2 := win.Rotate90()
	if r2 != 180 {
		t.Errorf("Expected 180 degrees, got %d", r2)
	}

	// Test flip
	f := win.ToggleFlipX()
	if !f {
		t.Errorf("Expected FlipX to be true")
	}
}

func TestPinManager(t *testing.T) {
	mgr := NewPinManager()
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))

	p1 := mgr.CreatePin(50, 50, 200, 200, img)
	p2 := mgr.CreatePin(100, 100, 200, 200, img)

	if mgr.Count() != 2 {
		t.Fatalf("Expected 2 pins, got %d", mgr.Count())
	}

	retrieved, exists := mgr.GetPin(p1.ID)
	if !exists || retrieved.ID != p1.ID {
		t.Errorf("Failed to retrieve pin %s", p1.ID)
	}

	mgr.ClosePin(p1.ID)
	if mgr.Count() != 1 {
		t.Errorf("Expected 1 pin after closing, got %d", mgr.Count())
	}

	mgr.ClearAll()
	if mgr.Count() != 0 {
		t.Errorf("Expected 0 pins after clear, got %d", mgr.Count())
	}
	_ = p2
}
