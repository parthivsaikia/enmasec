package utils

import (
	"os"
	"path/filepath"
)

func GetEnmasecDirLocation() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", ErrGetHomeDir(err)
	}
	vaultLocation := filepath.Join(homeDir, "enmasec")
	return vaultLocation, nil
}
