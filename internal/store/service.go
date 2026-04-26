package store

import (
	"os"
	"path/filepath"
)

func CreateService(vaultPath, serviceName string, indexMapData []byte) error {
	servicePath := filepath.Join(vaultPath, serviceName)
	if err := os.Mkdir(servicePath, 0o700); err != nil {
		return err
	}
	// write the given data back to the index file
	err := WriteIndexFile(vaultPath, indexMapData)
	if err != nil {
		return err
	}
	return nil
}
