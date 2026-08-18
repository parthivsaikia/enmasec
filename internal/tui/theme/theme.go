package theme

import (
	"fmt"
	"path/filepath"

	"github.com/parthivsaikia/enmasec/internal/config"
	"github.com/parthivsaikia/enmasec/internal/store"
	"gopkg.in/yaml.v3"
)

type Theme struct {
	Name string
	Pane pane
	List list
	Form form
}

type pane struct {
	Foreground  string `yaml:"fg"`
	Background  string `yaml:"bg"`
	BorderColor string `yaml:"border-color"`
}

type list struct {
	HighlightedForeground string `yaml:"highlighted-fg"`
	HighlightedBackground string `yaml:"highlighted-bg"`
	Foreground            string `yaml:"fg"`
	Background            string `yaml:"bg"`
}

type form struct {
	BorderColor      string `yaml:"border-color"`
	Foreground       string `yaml:"fg"`
	Background       string `yaml:"bg"`
	Input            string `yaml:"input"`
	HighlightedInput string `yaml:"highlighted-input"`
}

func GetCurrTheme() *Theme {
	theme := config.Config.Theme
	if theme == "" {
		return &ThemeForest
	}
	if isPreInstalledTheme(theme) {
		return PreInstalledThemes[theme]
	}
	themeFile := getThemeFile(theme)
	themeData, err := store.ReadFile(themeFile)
	if err != nil {
		return &ThemeAurora
	}
	var t Theme
	err = yaml.Unmarshal(themeData, &t)
	if err != nil {
		return &ThemeAurora
	}
	return &t
}

func isPreInstalledTheme(theme string) bool {
	_, ok := PreInstalledThemes[theme]
	return ok
}

func getThemeFile(theme string) string {
	configDir := store.GetEnmasecConfigDirLocation()
	themeFile := filepath.Join(configDir, "themes", fmt.Sprintf("%s.yaml", theme))
	return themeFile
}
