package theme

import (
	"charm.land/lipgloss/v2"
)

type Style struct {
	PaneStyles *lipgloss.Style
	ListStyles *listStyle
	FormStyles *lipgloss.Style
}

type listStyle struct {
	HighlightedStyles lipgloss.Style
	NormalStyles      lipgloss.Style
}

func buildStyle(theme *Theme) *Style {
	paneStyle := buildPaneStyle(theme)
	listStyle := buildListStyle(theme)
	return &Style{
		PaneStyles: &paneStyle,
		ListStyles: listStyle,
	}
}

func buildPaneStyle(theme *Theme) lipgloss.Style {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Pane.Foreground)).
		Background(lipgloss.Color(theme.Pane.Background)).
		BorderForeground(lipgloss.Color(theme.Pane.BorderColor)).
		BorderBottom(true).
		BorderTop(true).
		BorderRight(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderBottomForeground(lipgloss.Color(theme.Pane.BorderColor)).
		BorderStyle(lipgloss.ThickBorder())
	return style
}

func buildListStyle(theme *Theme) *listStyle {
	highlightedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.List.HighlightedForeground)).
		Background(lipgloss.Color(theme.List.HighlightedBackground))

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.List.Foreground)).
		Background(lipgloss.Color(theme.List.Background))
	return &listStyle{
		HighlightedStyles: highlightedStyle,
		NormalStyles:      normalStyle,
	}
}

func GetCurrentStyle() *Style {
	theme := GetCurrTheme()
	return buildStyle(theme)
}
