package store

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/parthivsaikia/enmasec/internal/encryption"
	"github.com/parthivsaikia/enmasec/internal/service"
)

func CreateService(path, password string) error {
	uuid := service.GenerateUUID()
	service.LocationUUIDMap.Put(uuid, path)
	// encrypt this map using age in the file vault/index.age
	mapData, err := json.Marshal(service.LocationUUIDMap)
	if err != nil {
		return err
	}
	encryptedData, err := encryption.EncryptAge(mapData, password)
	if err != nil {
		return err
	}
	fmt.Print(encryptedData)

	if err := os.Mkdir(path, 0o700); err != nil {
		return err
	}
	return nil
}
