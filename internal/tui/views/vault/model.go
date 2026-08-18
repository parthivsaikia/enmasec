package vault

import (
	"fmt"
	"log"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/parthivsaikia/enmasec/internal/models"
	"github.com/parthivsaikia/enmasec/internal/tui/components/input"
	"github.com/parthivsaikia/enmasec/internal/tui/components/list"
	"github.com/parthivsaikia/enmasec/internal/tui/keymaps"
	"github.com/parthivsaikia/enmasec/internal/tui/messages"
	"github.com/parthivsaikia/enmasec/internal/tui/theme"
	"golang.org/x/term"
)

type lockedStatus int

const (
	locked = iota
	unlocking
	unlocked
)

type entry struct {
	vault  *models.Vault
	status lockedStatus
}

type Model struct {
	VaultCreateForm *input.VaultCreateForm
	IsScreenOpen    bool
	searchBox       textinput.Model
	entries         []entry
	vaultList       *list.Model
	style           *lipgloss.Style
	width           int
	height          int
}

func New(vaults []*models.Vault) Model {
	var entries []entry
	for _, v := range vaults {
		e := entry{
			vault:  v,
			status: locked,
		}
		entries = append(entries, e)
	}
	log.Print("entries: ", generateVaultList(entries))
	return Model{
		entries:         entries,
		VaultCreateForm: input.NewVaultCreateForm(),
		IsScreenOpen:    false,
		vaultList:       list.New(generateVaultList(entries)),
		style:           theme.GetCurrentStyle().PaneStyles,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		m.vaultList, cmd = m.vaultList.Update(msg)
		cmds = append(cmds, cmd)

		cmd = m.VaultCreateForm.Update(msg)
		cmds = append(cmds, cmd)

		m.height = msg.Height
		m.width = msg.Width / 3

		return m, tea.Batch(cmds...)
	}

	if m.VaultCreateForm.IsOpen() {
		m.IsScreenOpen = true
		return m.UpdateCreateVaultInput(msg)
	}

	return m.UpdateList(msg)
}

func (m Model) UpdateList(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymaps.VaultKeys.OpenVaultCreateForm):
			cmd := m.VaultCreateForm.Open()
			cmds = append(cmds, cmd)
		default:
			var cmd tea.Cmd
			m.vaultList, cmd = m.vaultList.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

	case messages.VaultCreateMsg:
		m.VaultCreateForm = input.NewVaultCreateForm()
	}
	return m, tea.Batch(cmds...)
}

func (m Model) UpdateCreateVaultInput(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			m.VaultCreateForm.Close()
			m.IsScreenOpen = false
			return m, nil
		}
	}
	cmd := m.VaultCreateForm.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	width, height, _ := term.GetSize(0)

	backGroundStr := m.style.Height(m.height).Width(m.width).Render(m.vaultList.View())
	log.Print("backGroundStr: ", backGroundStr)

	if !m.VaultCreateForm.IsOpen() {
		return backGroundStr
	}

	overlayStr := m.VaultCreateForm.View()

	alignedStyle := lipgloss.NewStyle().Align(lipgloss.Center)
	alignedOverlay := alignedStyle.Render(overlayStr)

	layers := []*lipgloss.Layer{
		lipgloss.NewLayer(backGroundStr),
		// By using lipgloss.Center for alignment, adding the errors
		// will naturally expand the overlay downwards while keeping it centered!
		lipgloss.NewLayer(alignedOverlay).X(width / 2).Y(height / 2).Z(2),
	}

	compositedView := lipgloss.NewCompositor(layers...)
	return compositedView.Render()
}

func generateVaultList(entries []entry) []string {
	var s []string
	lockedIcon := "[L]"
	unlockedIcon := "[U]"
	for _, e := range entries {
		var lockedStatus string
		if e.status == locked {
			lockedStatus = lockedIcon
		} else {
			lockedStatus = unlockedIcon
		}
		line := fmt.Sprintf("%s %s", lockedStatus, e.vault.Name)
		s = append(s, line)
	}
	return s
}
