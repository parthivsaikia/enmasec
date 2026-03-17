package store

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/parthivsaikia/enmasec/internal/encryption"
	"github.com/parthivsaikia/enmasec/internal/utils"
)

func CreateVault(vaultLocation, password string, hash []byte) error {
	if err := os.MkdirAll(vaultLocation, 0o700); err != nil {
		return fmt.Errorf("unable to create vault %w", err)
	}
	keyFile := filepath.Join(vaultLocation, "key.age")
	file, err := os.Create(keyFile)
	if err != nil {
		return fmt.Errorf("unable to create key file %w", err)
	}
	defer file.Close()

	if _, err := file.Write(append(hash, []byte("\n")...)); err != nil {
		return err
	}

	data, err := encryption.EncryptAge([]byte(KEY_FILE_TEXT), password)
	if err != nil {
		return fmt.Errorf("unable to encrypt key file: %w", err)
	}

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("unable to write to file %s", file.Name())
	}
	return nil
}
