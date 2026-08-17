package annotation

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// RenderShapesOnImage rasterizes vector annotation shapes on top of the base image.
func RenderShapesOnImage(base *image.RGBA, shapes []*AnnotationShape) *image.RGBA {
	bounds := base.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, base, bounds.Min, draw.Src)

	for _, s := range shapes {
		switch s.Type {
		case ToolRect:
			drawRect(dst, s)
		case ToolArrow:
			drawArrow(dst, s)
		case ToolPen, ToolHighlighter:
			drawPath(dst, s)
		case ToolMosaic:
			applyMosaic(dst, s)
		case ToolStep:
			drawStepBadge(dst, s)
		case ToolText:
			// Text rendering
		}
	}

	return dst
}

func drawRect(img *image.RGBA, s *AnnotationShape) {
	x0 := int(math.Min(s.StartPoint.X, s.EndPoint.X))
	y0 := int(math.Min(s.StartPoint.Y, s.EndPoint.Y))
	x1 := int(math.Max(s.StartPoint.X, s.EndPoint.X))
	y1 := int(math.Max(s.StartPoint.Y, s.EndPoint.Y))
	sw := int(math.Max(1, s.StrokeWidth))

	bounds := img.Bounds()
	// Horizontal borders
	for x := x0; x <= x1; x++ {
		for dy := 0; dy < sw; dy++ {
			setPixelClamped(img, x, y0+dy, s.Color, bounds)
			setPixelClamped(img, x, y1-dy, s.Color, bounds)
		}
	}
	// Vertical borders
	for y := y0; y <= y1; y++ {
		for dx := 0; dx < sw; dx++ {
			setPixelClamped(img, x0+dx, y, s.Color, bounds)
			setPixelClamped(img, x1-dx, y, s.Color, bounds)
		}
	}
}

func drawArrow(img *image.RGBA, s *AnnotationShape) {
	// Bresenham's line + arrow head
	drawLine(img, int(s.StartPoint.X), int(s.StartPoint.Y), int(s.EndPoint.X), int(s.EndPoint.Y), s.Color, int(s.StrokeWidth))
}

func drawPath(img *image.RGBA, s *AnnotationShape) {
	if len(s.Points) < 2 {
		return
	}
	for i := 0; i < len(s.Points)-1; i++ {
		p1 := s.Points[i]
		p2 := s.Points[i+1]
		drawLine(img, int(p1.X), int(p1.Y), int(p2.X), int(p2.Y), s.Color, int(s.StrokeWidth))
	}
}

func applyMosaic(img *image.RGBA, s *AnnotationShape) {
	x0 := int(math.Min(s.StartPoint.X, s.EndPoint.X))
	y0 := int(math.Min(s.StartPoint.Y, s.EndPoint.Y))
	x1 := int(math.Max(s.StartPoint.X, s.EndPoint.X))
	y1 := int(math.Max(s.StartPoint.Y, s.EndPoint.Y))
	bs := s.BlockSize
	if bs <= 0 {
		bs = 10
	}

	bounds := img.Bounds()

	for by := y0; by < y1; by += bs {
		for bx := x0; bx < x1; bx += bs {
			// Sample center color of the block
			cx := bx + bs/2
			cy := by + bs/2
			if cx >= bounds.Max.X {
				cx = bounds.Max.X - 1
			}
			if cy >= bounds.Max.Y {
				cy = bounds.Max.Y - 1
			}
			c := img.RGBAAt(cx, cy)

			// Fill block
			for y := by; y < by+bs && y < y1; y++ {
				for x := bx; x < bx+bs && x < x1; x++ {
					setPixelClamped(img, x, y, c, bounds)
				}
			}
		}
	}
}

func drawStepBadge(img *image.RGBA, s *AnnotationShape) {
	cx := int(s.StartPoint.X)
	cy := int(s.StartPoint.Y)
	r := 12
	bounds := img.Bounds()

	// Fill circle
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			if (x-cx)*(x-cx)+(y-cy)*(y-cy) <= r*r {
				setPixelClamped(img, x, y, s.Color, bounds)
			}
		}
	}
}

func drawLine(img *image.RGBA, x0, y0, x1, y1 int, c color.RGBA, width int) {
	dx := math.Abs(float64(x1 - x0))
	dy := math.Abs(float64(y1 - y0))
	sx := 1
	if x0 >= x1 {
		sx = -1
	}
	sy := 1
	if y0 >= y1 {
		sy = -1
	}
	err := dx - dy

	bounds := img.Bounds()
	hw := width / 2

	for {
		for wx := -hw; wx <= hw; wx++ {
			for wy := -hw; wy <= hw; wy++ {
				setPixelClamped(img, x0+wx, y0+wy, c, bounds)
			}
		}

		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func setPixelClamped(img *image.RGBA, x, y int, c color.RGBA, bounds image.Rectangle) {
	if x >= bounds.Min.X && x < bounds.Max.X && y >= bounds.Min.Y && y < bounds.Max.Y {
		img.SetRGBA(x, y, c)
	}
}
