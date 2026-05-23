package components

import (
	"fmt"
	"strings"
	"syscall"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"golang.org/x/term"
)

var (
	purple = lipgloss.Color("99")
	gray   = lipgloss.Color("245")
	red    = lipgloss.Red

	headerStyle  = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
	cellStyle    = lipgloss.NewStyle().Padding(0, 1)
	oddRowStyle  = cellStyle.Foreground(gray)
	evenRowStyle = cellStyle.Foreground(red)
)

func PasswordPrompt(prompt string) (string, error) {
	// TODO: add a dynamic checklist of each of the password requirements which gets ticked once each requirement is met
	fmt.Printf("%s", prompt)
	bytes, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", err
	}
	fmt.Println()
	return strings.TrimSpace(string(bytes)), nil
}

func VaultTable(currentVaultRow int, rows [][]string) error {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		Headers("Name", "Location").
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == currentVaultRow:
				return evenRowStyle
			default:
				return oddRowStyle
			}
		}).
		Rows(rows...)
	if _, err := lipgloss.Println(t); err != nil {
		return err
	}
	return nil
}
