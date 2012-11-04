package entity

type EventType string

type Event struct {
	// Simple example: X damages Y and Z sees. Z is the Target,
	// Y is the Caller, and X is the Sender. EventType might be
	// something like "saw combat"

	// The target of the event (who the event is given to)
	Target EntityID

	// The caller of the event (who the event was triggered on)
	Caller EntityID

	// The sender of the event (who triggered the event)
	Sender EntityID

	EventType EventType
}

// Embed this to silently discard all events
type NoEvents struct{}

func (NoEvents) AcceptEvent(Event) {}
