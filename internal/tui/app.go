package tui

import (
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/parthivsaikia/enmasec/internal/tui/views/home"
)

func App() error {
	logfilePath := "log.txt"
	f, err := tea.LogToFile(logfilePath, "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	p := tea.NewProgram(home.InitialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return nil
}
