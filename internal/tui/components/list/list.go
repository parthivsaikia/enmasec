package list

import (
	"log"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/parthivsaikia/enmasec/internal/tui/keymaps"
	"github.com/parthivsaikia/enmasec/internal/tui/theme"
)

type Model struct {
	list         []string
	viewport     viewport.Model
	cursor       int
	focusedStyle lipgloss.Style
	normalStyle  lipgloss.Style
	height       int
	width        int
}

func New(list []string) *Model {
	m := Model{
		viewport:     viewport.Model{},
		cursor:       0,
		focusedStyle: theme.GetCurrentStyle().ListStyles.HighlightedStyles,
		normalStyle:  theme.GetCurrentStyle().ListStyles.NormalStyles,
		list:         list,
	}
	return &m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport = viewport.New(viewport.WithHeight(msg.Height), viewport.WithWidth(msg.Width/3))
		m.viewport.SetContent(m.generateStyledList())
		m.height = msg.Height
		m.width = msg.Width / 3
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keymaps.CommonKeyMap.Down):
			if m.cursor < len(m.list)-1 {
				m.cursor++
			}
		case key.Matches(msg, keymaps.CommonKeyMap.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		}
	}
	return &m, nil
}

func (m Model) View() string {
	content := m.generateStyledList()
	m.viewport.SetContent(content)
	return m.viewport.View()
}

func (m Model) generateStyledList() string {
	hasBorder := m.normalStyle.GetBorderTopSize() > 0
	var s strings.Builder
	for i, e := range m.list {
		var styledBlock string
		log.Print("element: ", e)
		if i == m.cursor {
			itemWidth := m.width - m.focusedStyle.GetHorizontalFrameSize()
			style := m.focusedStyle.Width(itemWidth)
			if !hasBorder {
				style = style.MarginBottom(1)
			}
			styledBlock = style.Render(e)
		} else {
			itemWidth := m.width - m.normalStyle.GetHorizontalFrameSize()
			style := m.normalStyle.Width(itemWidth)
			if !hasBorder {
				style = style.MarginBottom(1)
			}
			styledBlock = style.Render(e)
		}
		if i > 0 && !hasBorder {
			s.WriteString("\n")
		}
		s.WriteString(styledBlock)
	}
	return s.String()
}
