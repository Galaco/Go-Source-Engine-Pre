package messages

import (
	"github.com/galaco/Gource-Engine/engine/core/event"
	"github.com/galaco/Gource-Engine/engine/input/keyboard"
)

type KeyHeld struct {
	event.Message
	Key keyboard.Key
}