package vault

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/parthivsaikia/enmasec/internal/models"
	"github.com/parthivsaikia/enmasec/internal/tui/components/input"
	"golang.org/x/term"
)

type lockedStatus int

const (
	locked = iota
	unlocking
	unlocked
)

type styles struct {
	title        lipgloss.Style
	item         lipgloss.Style
	selectedItem lipgloss.Style
	pagination   lipgloss.Style
	help         lipgloss.Style
	quitText     lipgloss.Style
}

func newStyles(darkBG bool) styles {
	var s styles
	s.title = lipgloss.NewStyle().MarginLeft(2)
	s.item = lipgloss.NewStyle().PaddingLeft(4)
	s.selectedItem = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	s.pagination = list.DefaultStyles(darkBG).PaginationStyle.PaddingLeft(4)
	s.help = list.DefaultStyles(darkBG).HelpStyle.PaddingLeft(4).PaddingBottom(1)
	s.quitText = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	return s
}

type itemDelegate struct {
	styles *styles
}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(entry)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.vault.Name)

	fn := d.styles.item.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return d.styles.selectedItem.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type entry struct {
	vault  *models.Vault
	status lockedStatus
}

func (e entry) FilterValue() string { return e.vault.Name }

type Model struct {
	Vaults          []*models.Vault
	vaultsList      list.Model
	VaultCreateForm input.VaultCreateForm
	IsScreenOpen    bool
}

func New(vaults []*models.Vault) Model {
	l := generateVaultList(vaults)
	return Model{
		vaults:          vaults,
		vaultsList:      l,
		vaultCreateForm: *input.NewVaultCreateForm(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	if _, ok := msg.(tea.WindowSizeMsg); ok {
		var cmd tea.Cmd

		m.vaultsList, cmd = m.vaultsList.Update(msg)
		cmds = append(cmds, cmd)

		cmd = m.vaultCreateForm.Update(msg)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)
	}

	if m.VaultCreateForm.IsOpen() {
		m.IsScreenOpen = true
		return m.UpdateCreateVaultInput(msg)
	}

	return m.UpdateList(msg)
}

func (m Model) UpdateList(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "n":
			m.VaultCreateForm.Open()
			return m, m.VaultCreateForm.Init()
		}
	}
	var cmd tea.Cmd
	m.vaultsList, cmd = m.vaultsList.Update(msg)
	return m, cmd
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

	backGroundStr := m.vaultsList.View()

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
