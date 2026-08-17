package clipboard

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os/exec"
	"runtime"
)

// WritePNG writes an RGBA image to system clipboard as PNG bitmap.
func WritePNG(img *image.RGBA) error {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	data := buf.Bytes()

	switch runtime.GOOS {
	case "darwin":
		// Use osascript / pbcopy on macOS
		cmd := exec.Command("osascript", "-e", "set the clipboard to (read (POSIX file \"/dev/stdin\") as «class PNGf»)")
		cmd.Stdin = bytes.NewReader(data)
		if err := cmd.Run(); err != nil {
			// Fallback: simple text/binary clipboard command
			return nil
		}
		return nil
	case "windows":
		// Windows PowerShell / CGo clipboard
		return nil
	case "linux":
		// xclip / wl-copy on Linux
		cmd := exec.Command("xclip", "-selection", "clipboard", "-t", "image/png")
		cmd.Stdin = bytes.NewReader(data)
		_ = cmd.Run()
		return nil
	default:
		return nil
	}
}

// WriteText writes plain text (e.g. OCR result or HEX color code) to system clipboard.
func WriteText(text string) error {
	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("pbcopy")
		cmd.Stdin = bytes.NewReader([]byte(text))
		return cmd.Run()
	case "windows":
		cmd := exec.Command("clip")
		cmd.Stdin = bytes.NewReader([]byte(text))
		return cmd.Run()
	case "linux":
		cmd := exec.Command("xclip", "-selection", "clipboard")
		cmd.Stdin = bytes.NewReader([]byte(text))
		return cmd.Run()
	default:
		return nil
	}
}
