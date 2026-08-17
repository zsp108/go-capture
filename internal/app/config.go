package app

// Config stores application runtime preferences.
type Config struct {
	HotkeyCapture string  `json:"hotkey_capture"` // e.g. "F1" or "Ctrl+Shift+A"
	HotkeyPin     string  `json:"hotkey_pin"`     // e.g. "F3"
	DefaultColor  string  `json:"default_color"`  // e.g. "#EF4444"
	DefaultStroke float64 `json:"default_stroke"` // e.g. 2.0
	SavePath      string  `json:"save_path"`
	AutoCopy      bool    `json:"auto_copy"`
	SoundEnabled  bool    `json:"sound_enabled"`
}

// DefaultConfig returns the standard out-of-the-box configuration.
func DefaultConfig() *Config {
	return &Config{
		HotkeyCapture: "Ctrl+Shift+A",
		HotkeyPin:     "F3",
		DefaultColor:  "#EF4444",
		DefaultStroke: 2.0,
		SavePath:      "",
		AutoCopy:      true,
		SoundEnabled:  true,
	}
}
