package keymaps

import (
	"charm.land/bubbles/v2/key"
)

type KeyMap struct {
	Up   key.Binding
	Down key.Binding
}

var CommonKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
	),
}
