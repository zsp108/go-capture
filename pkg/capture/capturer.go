package capture

import (
	"image"
	"image/color"
	"image/draw"
)

// DisplayInfo represents display resolution and bounds.
type DisplayInfo struct {
	Index       int
	Bounds      image.Rectangle
	ScaleFactor float64 // For Retina / HiDPI (e.g., 2.0 on Retina Mac)
	IsPrimary   bool
}

// ScreenCapturer defines the interface for native 0-copy screen capture.
type ScreenCapturer interface {
	// GetDisplays returns all active monitors/displays.
	GetDisplays() ([]DisplayInfo, error)

	// CaptureDisplay captures the entire specified display.
	CaptureDisplay(displayIndex int) (*image.RGBA, error)

	// CaptureRect captures a specific sub-rectangle across the virtual screen.
	CaptureRect(rect image.Rectangle) (*image.RGBA, error)
}

// CropImage returns a sub-image cropped to the target rectangle.
func CropImage(src *image.RGBA, rect image.Rectangle) *image.RGBA {
	bounds := src.Bounds()
	intersect := rect.Intersect(bounds)
	if intersect.Empty() {
		return image.NewRGBA(image.Rect(0, 0, 1, 1))
	}

	dst := image.NewRGBA(image.Rect(0, 0, intersect.Dx(), intersect.Dy()))
	draw.Draw(dst, dst.Bounds(), src, intersect.Min, draw.Src)
	return dst
}

// GetPixelColor gets the RGBA color of a pixel at (x, y).
func GetPixelColor(img *image.RGBA, x, y int) color.RGBA {
	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return color.RGBA{0, 0, 0, 255}
	}
	return img.RGBAAt(x, y)
}
