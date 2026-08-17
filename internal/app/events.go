package app

// EventType defines application event types.
type EventType string

const (
	EventStartCapture   EventType = "start_capture"
	EventSelectionDone  EventType = "selection_done"
	EventColorSampled   EventType = "color_sampled"
	EventAnnotationDone EventType = "annotation_done"
	EventPinCreated     EventType = "pin_created"
	EventOCRRequested   EventType = "ocr_requested"
	EventCaptureSaved   EventType = "capture_saved"
	EventCaptureCanceled EventType = "capture_canceled"
)

// AppEvent represents a unified event payload across the architecture.
type AppEvent struct {
	Type    EventType   `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}
