package store

import (
	"os"
	"path/filepath"
)

func CreateService(vaultPath, serviceName string, indexMapData []byte) error {
	servicePath := filepath.Join(vaultPath, serviceName)
	indexFilePath := filepath.Join(vaultPath, "index.age")
	if err := os.Mkdir(servicePath, 0o700); err != nil {
		return err
	}
	// write the given data back to the index file
	err := WriteFile(indexMapData, indexFilePath)
	if err != nil {
		os.Remove(servicePath)
		return err
	}
	return nil
}
