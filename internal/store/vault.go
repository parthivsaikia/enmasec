package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/parthivsaikia/enmasec/internal/encryption"
	"github.com/parthivsaikia/enmasec/internal/utils"
)

func CreateVault(vaultLocation, password string) error {
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

	data, err := encryption.EncryptAge([]byte(KEY_FILE_TEXT), password)
	if err != nil {
		return fmt.Errorf("unable to encrypt key file: %w", err)
	}

	if _, err := kf.Write(data); err != nil {
		return fmt.Errorf("unable to write to file %s", kf.Name())
	}
	return nil
}

func Unlock(vaultPath, password string) error {
	keyfile := filepath.Join(vaultPath, "key.age")
	if !utils.CheckFileExists(keyfile) {
		fmt.Println(keyfile)
		return fmt.Errorf("key file doesn't exist")
	}
	content, _ := os.ReadFile(keyfile)

	key, err := encryption.DecryptAge(password, content)
	if err != nil {
		return fmt.Errorf("unable to decrypt file: %w", err)
	}

	if string(key) != (KEY_FILE_TEXT) {
		return fmt.Errorf("key text doesn't match")
	}
	return nil
}
