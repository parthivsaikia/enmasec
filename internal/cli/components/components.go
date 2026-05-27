package components

import (
	"fmt"
	"strings"
	"syscall"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/parthivsaikia/enmasec/internal/models"
	"github.com/parthivsaikia/enmasec/internal/validation"
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

func AccountCreationREPL(account *models.Account) error {
	var more bool
	coreForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter your username").
				Validate(validation.ValidateAccountName).
				Value(&(account.Username)),
			huh.NewInput().
				Title("Enter your password").
				Value(&(account.Password)),
			huh.NewConfirm().
				Title("Do you have more fields to enter").
				Value(&more),
		),
	)
	if err := coreForm.Run(); err != nil {
		return err
	}
	for more {
		var key, val string
		metaDataForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Enter key").
					Validate(validation.ValidateAccountMetaDataKey).
					Value(&key),
				huh.NewInput().
					Title("Enter value").
					Validate(validation.ValidateAccountMetaDataKey).
					Value(&val),
				huh.NewConfirm().
					Title("Do you have more fields to enter").
					Value(&more),
			),
		)
		if err := metaDataForm.Run(); err != nil {
			return err
		}
		account.Metadata[key] = val
	}
	return nil
}
