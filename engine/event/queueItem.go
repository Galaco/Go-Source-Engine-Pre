package event

import (
	"github.com/galaco/go-me-engine/engine/interfaces"
	"github.com/galaco/go-me-engine/engine/core"
)

// Event Queue item.
// Contains the event name, and a message
type QueueItem struct {
	EventName core.EventId
	Message interfaces.IMessage
}
