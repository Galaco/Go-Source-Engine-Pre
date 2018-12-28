package input

import (
	"github.com/galaco/Gource-Engine/client/input/keyboard"
	"github.com/galaco/Gource-Engine/client/messages"
	"github.com/galaco/Gource-Engine/core/event"
)

// Keyboard key wrapper
type Keyboard struct {
	keysDown [1024]bool
}

// IsKeyDown Check if a specific key is pressed
func (keyboard *Keyboard) IsKeyDown(key keyboard.Key) bool {
	return keyboard.keysDown[int(key)]
}

// ReceiveMessage Event manager message receiver.
// Used to catch key events from the window library
func (keyboard *Keyboard) ReceiveMessage(message event.IMessage) {
	switch message.GetType() {
	case messages.TypeKeyDown:
		keyboard.keysDown[int(message.(*messages.KeyDown).Key)] = true
	case messages.TypeKeyReleased:
		keyboard.keysDown[int(message.(*messages.KeyReleased).Key)] = false
	}
}

func (keyboard *Keyboard) SendMessage() event.IMessage {
	return nil
}

var staticKeyboard Keyboard

// GetKeyboard return static keyboard
func GetKeyboard() *Keyboard {
	return &staticKeyboard
}