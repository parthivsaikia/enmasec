package theme

import (
	"charm.land/lipgloss/v2"
)

type Style struct {
	PaneStyles *lipgloss.Style
	ListStyles *lipgloss.Style
	FormStyles *lipgloss.Style
}

func buildStyle(theme *Theme) *Style {
	paneStyle := buildPaneStyle(theme)
	return &Style{
		PaneStyles: &paneStyle,
	}
}

func buildPaneStyle(theme *Theme) lipgloss.Style {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Pane.Foreground)).
		Background(lipgloss.Color(theme.Pane.Background)).
		BorderForeground(lipgloss.Color(theme.Pane.BorderColor))

	if theme.Pane.BorderRadius == "rounded" {
		style = style.Border(lipgloss.RoundedBorder())
	}

	return style
}

func GetCurrentStyle() *Style {
	theme := GetCurrTheme()
	return buildStyle(theme)
}
