package keymaps

import (
	"charm.land/bubbles/v2/key"
)

type VaultKeyMap struct {
	OpenVaultCreateForm key.Binding
}

var VaultKeys = VaultKeyMap{
	OpenVaultCreateForm: key.NewBinding(
		key.WithKeys("n"),
	),
}
