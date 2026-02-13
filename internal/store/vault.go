package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/parthivsaikia/enmasec/internal/encryption"
	"github.com/parthivsaikia/enmasec/internal/model"
)

func (s *Store) CreateVault(vault model.Vault, masterPassword string) error {
	if _, err := os.Stat(s.VaultPath); !os.IsNotExist(err) {
		return fmt.Errorf("vault %s already exixts", vault.Name)
	}
	if err := os.Mkdir(s.VaultPath, 0o755); err != nil {
		return fmt.Errorf("unable to create vault %v", err)
	}
	keyFile := filepath.Join(s.VaultPath, vault.Name)
	file, err := os.Create(keyFile)
	if err != nil {
		return fmt.Errorf("unable to create key file %v", err)
	}
	defer file.Close()
	data, err := encryption.Encrypt([]byte(KEY_FILE_TEXT), masterPassword)
	if err != nil {
		return fmt.Errorf("unable to encrypt key file: %v", err)
	}
	file.Write(data)
	return nil
}
