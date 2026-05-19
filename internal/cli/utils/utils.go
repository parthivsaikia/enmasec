package utils

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/parthivsaikia/enmasec/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
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

func ResolveVault(cmd *cobra.Command) (string, error) {
	vault, err := cmd.Flags().GetString("vault")
	if err != nil {
		return "", err
	}
	if vault == "" {
		vault = config.Config.CurrentVault
	}
	return vault, err
}
