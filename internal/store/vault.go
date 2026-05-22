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
