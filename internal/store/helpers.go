package store

import (
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
