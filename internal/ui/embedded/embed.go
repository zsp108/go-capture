package embedded

import (
	"embed"
	"io/fs"
)

//go:embed index.html style.css app.js
var Assets embed.FS

// GetSubFS returns the embedded assets filesystem.
func GetSubFS() fs.FS {
	return Assets
}
