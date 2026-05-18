package components

import (
	"fmt"
	"strings"
	"syscall"

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
