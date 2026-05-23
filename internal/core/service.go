package core

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/parthivsaikia/enmasec/internal/config"
	"github.com/parthivsaikia/enmasec/internal/encryption"
	"github.com/parthivsaikia/enmasec/internal/models"
	"github.com/parthivsaikia/enmasec/internal/store"
)

func CreateService(vaultName, serviceName, key string) (*models.VaultIndex, *models.RuntimeIndex, error) {
	vaultPath := config.Config.Vaults[vaultName]
	id := uuid.New()

	indexData, err := DecryptVaultIndex(vaultName, string(key))
	if err != nil {
		return nil, nil, err
	}

	vaultIndex, runtimeIndex, err := OpenVaultIndex(indexData)
	if err != nil {
		return nil, nil, err
	}
	if _, ok := runtimeIndex.ServiceNameToID[serviceName]; ok {
		return nil, nil, fmt.Errorf("service %s already exists", serviceName)
	}

	vaultIndex.Services[id] = &models.ServiceEntry{
		Name:     serviceName,
		Accounts: make(map[uuid.UUID]*models.AccountEntry),
	}

	runtimeIndex.ServiceNameToID[serviceName] = id
	runtimeIndex.ServiceIDToName[id] = serviceName
	runtimeIndex.AccountNameToID[id] = make(map[string]uuid.UUID)
	runtimeIndex.AccountIDToName[id] = make(map[uuid.UUID]string)

	indexBytes, err := json.Marshal(vaultIndex)
	if err != nil {
		return nil, nil, err
	}

	encryptedIndexMapData, err := encryption.EncryptAge(indexBytes, key)
	if err != nil {
		return nil, nil, err
	}
	err = store.CreateService(vaultPath, id.String(), encryptedIndexMapData)
	if err != nil {
		delete(vaultIndex.Services, id)
		delete(runtimeIndex.ServiceNameToID, serviceName)
		delete(runtimeIndex.ServiceIDToName, id)
		delete(runtimeIndex.AccountNameToID, id)
		delete(runtimeIndex.AccountIDToName, id)
		return nil, nil, err
	}
	return vaultIndex, runtimeIndex, nil
}

func BuildRuntimeIndex(v *models.VaultIndex) *models.RuntimeIndex {
	r := &models.RuntimeIndex{
		ServiceNameToID: make(map[string]uuid.UUID),
		ServiceIDToName: make(map[uuid.UUID]string),
		AccountNameToID: make(map[uuid.UUID]map[string]uuid.UUID),
		AccountIDToName: make(map[uuid.UUID]map[uuid.UUID]string),
	}

	for svcID, svc := range v.Services {
		r.ServiceNameToID[svc.Name] = svcID
		r.ServiceIDToName[svcID] = svc.Name
		r.AccountNameToID[svcID] = make(map[string]uuid.UUID)
		r.AccountIDToName[svcID] = make(map[uuid.UUID]string)

		for acctID, acct := range svc.Accounts {
			r.AccountNameToID[svcID][acct.Name] = acctID
			r.AccountIDToName[svcID][acctID] = acct.Name
		}
	}

	return r
}

func OpenVaultIndex(indexData []byte) (*models.VaultIndex, *models.RuntimeIndex, error) {
	var v models.VaultIndex
	if err := json.Unmarshal(indexData, &v); err != nil {
		return nil, nil, err
	}
	r := BuildRuntimeIndex(&v)
	return &v, r, nil
}

func ListService(vaultName, key string) error {
	indexData, err := DecryptVaultIndex(vaultName, key)
	if err != nil {
		return err
	}
	vi, _, err := OpenVaultIndex(indexData)
	if err != nil {
		return err
	}
	for _, svc := range vi.Services {
		fmt.Printf("%s\n", svc.Name)
	}
	return nil
}
