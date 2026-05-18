package core

import (
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/parthivsaikia/enmasec/internal/config"
	"github.com/parthivsaikia/enmasec/internal/encryption"
	"github.com/parthivsaikia/enmasec/internal/models"
	"github.com/parthivsaikia/enmasec/internal/store"
)

var LocationUUIDMap = models.BiMap{
	ForwardMap: map[uuid.UUID]string{},
	ReverseMap: map[string]uuid.UUID{},
}

func CreateVault(dir, vaultName, password string) error {
	vaultLocation := filepath.Join(dir, vaultName)
	secretKey := encryption.RandomByte(32)
	encryptedKey, err := encryption.EncryptAge(secretKey, password)
	if err != nil {
		return fmt.Errorf("unable to encrypt: %w", err)
	}
	err = store.CreateVaultStore(vaultLocation, password, encryptedKey)
	if err != nil {
		return err
	}
	config.Config.CurrentVault = vaultName
	config.Config.Vaults[vaultName] = vaultLocation
	if err := config.Save(); err != nil {
		return fmt.Errorf("couldn't save config: %w", err)
	}
	return nil
}

func UnlockVault(vaultName, password string) ([]byte, error) {
	vaultLocation := config.Config.Vaults[vaultName]
	keyFile := filepath.Join(vaultLocation, "key.age")
	content, err := store.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	data, err := encryption.DecryptAge(password, content)
	if err != nil {
		return nil, err
	}
	return data, nil
}
