package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/parthivsaikia/enmasec/internal/config"
	"github.com/parthivsaikia/enmasec/internal/models"
	"github.com/parthivsaikia/enmasec/internal/store"
)

var LocationUUIDMap = models.BiMap{
	ForwardMap: map[uuid.UUID]string{},
	ReverseMap: map[string]uuid.UUID{},
}

func InitIndexMap() {
	biMapData, err := store.DecryptIndexFile(vault, password)
}

func CreateVault(vaultLocation, password, vaultName string) error {
	err := store.CreateVault(vaultLocation, password)
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
