package messages

import (
	tea "charm.land/bubbletea/v2"
	"github.com/parthivsaikia/enmasec/internal/core"
)

func VaultCreatedMsg(vaultPath, vaultName, password string) tea.Msg {
	err := core.CreateVault(vaultPath, vaultName, password)
	if err != nil {
		return err
	}
	return nil
}
