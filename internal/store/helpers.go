package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

func GetEnmasecDirLocation() string {
	dataDir := xdg.DataHome
	enmasecDirLocation := filepath.Join(dataDir, "enmasec")
	return enmasecDirLocation
}

func GetEnmasecConfigDirLocation() string {
	configDir := xdg.ConfigHome
	enmasecConfigDirLocation := filepath.Join(configDir, "enmasec")
	return enmasecConfigDirLocation
}

func CheckFileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func ReadFile(filePath string) ([]byte, error) {
	if !CheckFileExists(filePath) {
		return nil, fmt.Errorf("%s file doesn't exist", filePath)
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %s: %w", filePath, err)
	}
	return content, nil
}

func WriteFile(data []byte, filePath string) error {
	if !CheckFileExists(filePath) {
		return fmt.Errorf("%s file doesn't exist", filePath)
	}
	err := os.WriteFile(filePath, data, 0o600)
	if err != nil {
		return fmt.Errorf("unable to write to file %s", filePath)
	}
	return nil
}

func DeleteFile(filePath string) error {
	if !CheckFileExists(filePath) {
		return fmt.Errorf("%s file doesn't exist", filePath)
	}
	err := os.RemoveAll(filePath)
	if err != nil {
		return err
	}
	return nil
}

func RenameFile(oldPath, newPath string) error {
	if !CheckFileExists(oldPath) {
		return fmt.Errorf("%s path doesn't exist", oldPath)
	}
	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}
	return nil
}
