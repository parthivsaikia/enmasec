package vault

import (
	"charm.land/bubbles/v2/list"
	"github.com/parthivsaikia/enmasec/internal/config"
	"github.com/parthivsaikia/enmasec/internal/models"
)

func generateVaultList(vaults []*models.Vault) list.Model {
	var listItems []list.Item
	var currentVaultIndex int

	for i, vault := range vaults {
		e := entry{
			vault:  vault,
			status: locked,
		}
		listItems = append(listItems, e)
		if vault.Name == config.Config.CurrentVault {
			currentVaultIndex = i
		}
	}

	s := newStyles(true)
	delegate := itemDelegate{styles: &s}

	l := list.New(listItems, delegate, 20, 20)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.Select(currentVaultIndex)
	l.DisableQuitKeybindings()
	return l
}
