package core

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/parthivsaikia/enmasec/internal/config"
	"github.com/parthivsaikia/enmasec/internal/encryption"
	"github.com/parthivsaikia/enmasec/internal/models"
	"github.com/parthivsaikia/enmasec/internal/store"
	"github.com/parthivsaikia/enmasec/internal/validation"
)

var LocationUUIDMap = models.BiMap{
	ForwardMap: map[uuid.UUID]string{},
	ReverseMap: map[string]uuid.UUID{},
}

func InitIndexMap() {
	validation.
}

func CreateVaultHelper(vaultLocation, password, vaultName string) error {
	secretKey := encryption.RandomByte(32)
	encryptedKey, err := encryption.EncryptAge(secretKey, password)
	if err != nil {
		return fmt.Errorf("unable to encrypt key file: %w", err)
	}

	err = store.CreateVault(vaultLocation, password, encryptedKey)
	if err != nil {
		return fmt.Errorf("couldn't create vault: %w", err)
	}
	config.Config.CurrentVault = vaultName
	config.Config.Vaults[vaultName] = vaultLocation
	if err := config.Save(); err != nil {
		return fmt.Errorf("couldn't save config: %w", err)
	}
	return nil
}
