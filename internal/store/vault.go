package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/parthivsaikia/enmasec/internal/encryption"
)

func CreateVaultStore(vaultLocation, password string, encryptedKey []byte) error {
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
	return nil
}

func Unlock(vaultPath, password string) ([]byte, error) {
	keyfile := filepath.Join(vaultPath, "key.age")
	if !CheckFileExists(keyfile) {
		fmt.Println(keyfile)
		return nil, fmt.Errorf("key file doesn't exist")
	}
	content, err := os.ReadFile(keyfile)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %s: %w", keyfile, err)
	}

	key, err := encryption.DecryptAge(password, content)
	if err != nil {
		return nil, fmt.Errorf("unable to decrypt file: %w", err)
	}
	return key, nil
}

func DecryptIndexFile(vaultPath string, password string) ([]byte, error) {
	indexFile := filepath.Join(vaultPath, "index.age")
	if !CheckFileExists(indexFile) {
		return nil, fmt.Errorf("index file doesn't exist")
	}
	content, err := os.ReadFile(indexFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %s: %w", indexFile, err)
	}
	data, err := encryption.DecryptAge(password, content)
	if err != nil {
		return nil, fmt.Errorf("unable to decrypt map content: %w", err)
	}
	return data, nil
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
