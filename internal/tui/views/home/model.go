package home

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/parthivsaikia/enmasec/internal/core"
	"github.com/parthivsaikia/enmasec/internal/models"
	"github.com/parthivsaikia/enmasec/internal/tui/messages"
	"github.com/parthivsaikia/enmasec/internal/tui/views/vault"
)

type currentPane int

const (
	vaultPane = iota
	servicePane
	accountPane
)

type model struct {
	// data
	vaults       []*models.Vault
	runtimeIndex *models.RuntimeIndex
	// panes
	vaultModel vault.Model
	// cursor
	pane currentPane
}

func InitialModel() model {
	vaults := core.GetVaults()
	return model{
		vaults:     vaults,
		vaultModel: vault.New(vaults),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "tab":
			m.pane = (m.pane + 1) % 3
		}
	}
	switch m.pane {
	case 0:
		m.vaultModel, cmd = m.vaultModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	v := tea.NewView(m.vaultModel.View())
	v.AltScreen = true
	return v
}
