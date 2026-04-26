package core

import (
	"encoding/json"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/parthivsaikia/enmasec/internal/encryption"
	"github.com/parthivsaikia/enmasec/internal/store"
)

func ServiceCreationHelper(vaultPath, serviceName, password string) error {
	// create the service path by joining vault path and service name
	servicePath := filepath.Join(vaultPath, serviceName)
	uuid := uuid.New()
	LocationUUIDMap.Put(uuid, servicePath)
	indexMapData, err := json.Marshal(LocationUUIDMap)
	if err != nil {
		return err
	}
	encryptedIndexMapData, err := encryption.EncryptAge(indexMapData, password)
	if err != nil {
		return err
	}
	err = store.CreateService(vaultPath, serviceName, encryptedIndexMapData)
	if err != nil {
		return err
	}
	return nil
}
