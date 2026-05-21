package store

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateVaultStore(vaultLocation, password string, encryptedKey, encryptedIndexByte []byte) error {
	if err := os.MkdirAll(vaultLocation, 0o700); err != nil {
		return fmt.Errorf("unable to create vault %w", err)
	}
	keyFile := filepath.Join(vaultLocation, "key.age")
	kf, err := os.Create(keyFile)
	if err != nil {
		return fmt.Errorf("unable to create key file %w", err)
	}
	defer kf.Close()

	indexFile := filepath.Join(vaultLocation, "index.age")
	iFile, err := os.Create(indexFile)
	if err != nil {
		return fmt.Errorf("unable to create index file %w", err)
	}
	defer iFile.Close()
	if _, err := kf.Write(encryptedKey); err != nil {
		return fmt.Errorf("unable to write to file %s", kf.Name())
	}

	if _, err := iFile.Write(encryptedIndexByte); err != nil {
		return fmt.Errorf("unable to write to file %s", iFile.Name())
	}
	return nil
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

func WriteIndexFile(vaultPath string, data []byte) error {
	indexFile := filepath.Join(vaultPath, "index.age")
	if !CheckFileExists(indexFile) {
		return fmt.Errorf("index file doesn't exist")
	}
	err := os.WriteFile(indexFile, data, 0o700)
	if err != nil {
		return fmt.Errorf("couldn't write to index file")
	}
	return nil
}
