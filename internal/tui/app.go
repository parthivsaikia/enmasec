package tui

import (
	"github.com/parthivsaikia/enmasec/internal/tui/views/home"
)

func App() error {
	if err := home.Home(); err != nil {
		return err
	}
	return nil
}
