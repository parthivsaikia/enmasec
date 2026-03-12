package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/parthivsaikia/enmasec/internal/encryption"
)

func CreateVault(vaultLocation, masterPassword string) error {
	if err := os.MkdirAll(vaultLocation, 0o700); err != nil {
		return fmt.Errorf("unable to create vault %w", err)
	}
	keyFile := filepath.Join(vaultLocation, "key.age")
	file, err := os.Create(keyFile)
	if err != nil {
		return fmt.Errorf("unable to create key file %w", err)
	}
	defer file.Close()
	data, err := encryption.Encrypt([]byte(KEY_FILE_TEXT), masterPassword)
	if err != nil {
		return fmt.Errorf("unable to encrypt key file: %w", err)
	}

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("unable to write to file %s", file.Name())
	}
	return nil
}
