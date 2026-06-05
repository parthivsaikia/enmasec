package home

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type screen int

const (
	vaultScreen screen = iota
	serviceScreen
	accountScreen
)

var (
	focusedBlockStyle = lipgloss.NewStyle().
				Width(20).
				Height(10).
				Border(lipgloss.BlockBorder(), true)

	unfocusedBlockStyle = focusedBlockStyle.
				Background(lipgloss.Color("#f1f1f1"))
)

type model struct {
	focused screen
}

func New() model {
	return model{
		focused: vaultScreen,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "tab":
			m.focused = nextScreen(m.focused)
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	var s strings.Builder

	s.WriteString(
		lipgloss.JoinHorizontal(
			lipgloss.Bottom,
			m.renderBlock(vaultScreen, "Vault"),
			m.renderBlock(serviceScreen, "Services"),
			m.renderBlock(accountScreen, "Account"),
		),
	)
	v := tea.NewView(s.String())
	v.AltScreen = true
	return v
}

func (m model) renderBlock(target screen, title string) string {
	if m.focused == target {
		return focusedBlockStyle.Render(title)
	}

	return unfocusedBlockStyle.Render(title)
}

func nextScreen(current screen) screen {
	return (current + 1) % 3
}

func Home() error {
	p := tea.NewProgram(New())

	_, err := p.Run()
	return err
}
