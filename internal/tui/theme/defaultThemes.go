package theme

var ThemeMidnight = Theme{
	Name: "Midnight",
	Pane: pane{Foreground: "#cdd6f4", Background: "#1e1e2e", BorderColor: "#89b4fa", BorderRadius: "1"},
	List: list{Foreground: "#a6adc8", Background: "#1e1e2e", HighlightedForeground: "#89b4fa", HighlightedBackground: "#313244", BorderRadius: "1"},
	Form: form{Foreground: "#cdd6f4", Background: "#1e1e2e", BorderColor: "#585b70", Input: "#181825", HighlightedInput: "#cba6f7", BorderRadius: "1"},
}

var ThemeForest = Theme{
	Name: "Forest",
	Pane: pane{Foreground: "#d3c6aa", Background: "#2b3339", BorderColor: "#a7c080", BorderRadius: "1"},
	List: list{Foreground: "#9da9a0", Background: "#2b3339", HighlightedForeground: "#a7c080", HighlightedBackground: "#323d43", BorderRadius: "1"},
	Form: form{Foreground: "#d3c6aa", Background: "#2b3339", BorderColor: "#4b565c", Input: "#232a2e", HighlightedInput: "#dbbc7f", BorderRadius: "1"},
}

var ThemeAurora = Theme{
	Name: "Aurora",
	Pane: pane{Foreground: "#edf2f4", Background: "#22223b", BorderColor: "#00f5d4", BorderRadius: "1"},
	List: list{Foreground: "#8d99ae", Background: "#22223b", HighlightedForeground: "#f15bb5", HighlightedBackground: "#4a4e69", BorderRadius: "1"},
	Form: form{Foreground: "#edf2f4", Background: "#22223b", BorderColor: "#9b5de5", Input: "#14142b", HighlightedInput: "#fee440", BorderRadius: "1"},
}

var ThemeSlate = Theme{
	Name: "Slate",
	Pane: pane{Foreground: "#c8c8c8", Background: "#21252b", BorderColor: "#5c6370", BorderRadius: "0"},
	List: list{Foreground: "#abb2bf", Background: "#21252b", HighlightedForeground: "#ffffff", HighlightedBackground: "#2c313a", BorderRadius: "0"},
	Form: form{Foreground: "#c8c8c8", Background: "#21252b", BorderColor: "#3b4048", Input: "#181a1f", HighlightedInput: "#61afef", BorderRadius: "0"},
}

var ThemeIvory = Theme{
	Name: "Ivory",
	Pane: pane{Foreground: "#657b83", Background: "#fdf6e3", BorderColor: "#b58900", BorderRadius: "1"},
	List: list{Foreground: "#839496", Background: "#fdf6e3", HighlightedForeground: "#268bd2", HighlightedBackground: "#eee8d5", BorderRadius: "1"},
	Form: form{Foreground: "#586e75", Background: "#fdf6e3", BorderColor: "#93a1a1", Input: "#eee8d5", HighlightedInput: "#cb4b16", BorderRadius: "1"},
}

var PreInstalledThemes = map[string]*Theme{
	"midnight": &ThemeMidnight,
	"forest":   &ThemeForest,
	"aurora":   &ThemeAurora,
	"slate":    &ThemeSlate,
	"ivory":    &ThemeIvory,
}
