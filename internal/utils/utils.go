package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/adrg/xdg"
	"golang.org/x/term"
)

func GetEnmasecDirLocation() string {
	dataDir := xdg.DataHome
	vaultLocation := filepath.Join(dataDir, "enmasec")
	return vaultLocation
}

func GetEnmasecConfigDir() string {
	configDir := xdg.ConfigHome
	enmaConfigDir := filepath.Join(configDir, "enmasec")
	return enmaConfigDir
}

func PasswordPrompt(prompt string) (string, error) {
	fmt.Printf("%s", prompt)
	bytes, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", err
	}
	fmt.Println()
	return strings.TrimSpace(string(bytes)), nil
}

func CheckFileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
