package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/term"
)

func GetEnmasecDirLocation() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", ErrGetHomeDir(err)
	}
	vaultLocation := filepath.Join(homeDir, ".enmasec")
	return vaultLocation, nil
}

func PasswordPrompt(prompt string) (string, error) {
	fmt.Printf("%s", prompt)
	bytes, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", err
	}
	fmt.Println()
	return string(bytes), nil
}
