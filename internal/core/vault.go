package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/google/uuid"
	"github.com/parthivsaikia/enmasec/internal/encryption"
	"github.com/parthivsaikia/enmasec/internal/models"
	"github.com/parthivsaikia/enmasec/internal/registry"
	"github.com/parthivsaikia/enmasec/internal/state"
	"github.com/parthivsaikia/enmasec/internal/store"
	"github.com/parthivsaikia/enmasec/internal/validation"
)

func CreateVault(dir, vaultName, password string) (*models.Vault, error) {
	vaultPath := filepath.Join(dir, vaultName)
	secretKey := encryption.RandomByte(32)
	encryptedKey, err := encryption.EncryptAge(secretKey, password)
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt: %w", err)
	}
	emptyIndex := &models.VaultIndex{
		Version:  1,
		Services: make(map[uuid.UUID]*models.ServiceEntry),
	}
	indexByte, err := json.Marshal(emptyIndex)
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt index data: %w", err)
	}
	encryptedIndexByte, err := encryption.EncryptAge(indexByte, string(secretKey))
	if err != nil {
		return nil, err
	}
	err = store.CreateVaultStore(vaultPath, password, encryptedKey, encryptedIndexByte)
	if err != nil {
		return nil, err
	}
	state.State.CurrentVault = vaultName
	registry.Registry.Vaults[vaultName] = vaultPath
	if err := state.Save(); err != nil {
		return nil, fmt.Errorf("couldn't save current state: %w", err)
	}
	if err := registry.Save(); err != nil {
		return nil, fmt.Errorf("couldn't save data to registry: %w", err)
	}

	return &models.Vault{
		Path: vaultPath,
		Name: vaultName,
	}, nil
}

func UnlockVault(vaultName, password string) ([]byte, error) {
	vaultPath := registry.Registry.Vaults[vaultName]
	keyFile := filepath.Join(vaultPath, "key.age")
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

func CheckoutVault(vaultName string) error {
	if _, ok := registry.Registry.Vaults[vaultName]; !ok {
		return fmt.Errorf("vault doesn't exist")
	} else {
		state.State.CurrentVault = vaultName
		if err := state.Save(); err != nil {
			return fmt.Errorf("unable to save current state: %w", err)
		}
	}
	return nil
}

func GetVaults() []*models.Vault {
	var vaults []*models.Vault
	for vaultName, vaultPath := range registry.Registry.Vaults {
		var v models.Vault
		v.Name = vaultName
		v.Path = vaultPath
		vaults = append(vaults, &v)
	}
	sort.Slice(vaults, func(i, j int) bool {
		return vaults[i].Name < vaults[j].Name
	})
	return vaults
}

func UpdateVault(vaultName, newVaultName, newDir, newPassword string, key []byte) error {
	vaultPath := registry.Registry.Vaults[vaultName]
	if newDir == "" {
		newDir = filepath.Dir(vaultPath)
	}

	if newVaultName == "" {
		newVaultName = vaultName
	}

	newVaultLocation := filepath.Join(newDir, newVaultName)

	if newVaultLocation != vaultPath {
		if err := os.Rename(vaultPath, newVaultLocation); err != nil {
			return err
		}
	}

	if newPassword != "" {
		if err := validation.CheckPasswordValid(newPassword); err != nil {
			return err
		}
		f := filepath.Join(newVaultLocation, "key.age")
		data, err := encryption.EncryptAge(key, newPassword)
		if err != nil {
			return err
		}
		err = os.WriteFile(f, data, 0o666)
		if err != nil {
			return err
		}
	}

	registry.Registry.Vaults[newVaultName] = newVaultLocation
	state.State.CurrentVault = newVaultName
	if newVaultName != vaultName {
		delete(registry.Registry.Vaults, vaultName)
	}
	if err := state.Save(); err != nil {
		return err
	}
	if err := registry.Save(); err != nil {
		return err
	}

	return nil
}

func DecryptVaultIndex(vaultName, key string) ([]byte, error) {
	vaultPath := registry.Registry.Vaults[vaultName]
	indexFilePath := filepath.Join(vaultPath, "index.age")
	indexBytes, err := store.ReadFile(indexFilePath)
	if err != nil {
		return nil, err
	}
	indexData, err := encryption.DecryptAge(key, indexBytes)
	if err != nil {
		return nil, err
	}
	return indexData, nil
}

func RepairVaultIndex(vaultName, key string) (*models.VaultIndex, *models.RuntimeIndex, error) {
	vaultPath := registry.Registry.Vaults[vaultName]
	indexFilePath := filepath.Join(vaultPath, "index.age")
	indexData, err := DecryptVaultIndex(vaultName, string(key))
	if err != nil {
		return nil, nil, err
	}

	vi, rt, err := OpenVaultIndex(indexData)
	if err != nil {
		return nil, nil, err
	}
	// repair services
	for id, service := range vi.Services {
		servicePath := filepath.Join(vaultPath, id.String())
		if !store.CheckFileExists(servicePath) {
			delete(vi.Services, id)
			delete(rt.ServiceIDToName, id)
			delete(rt.ServiceNameToID, service.Name)
		} else {
			for acctID := range service.Accounts {
				acctPath := filepath.Join(servicePath, acctID.String()+".age")
				if !store.CheckFileExists(acctPath) {
					acctName := vi.Services[id].Accounts[acctID].Name
					delete(vi.Services[id].Accounts, acctID)
					delete(rt.AccountNameToID[id], acctName)
					delete(rt.AccountIDToName[id], acctID)
				}
			}
		}

	}

	indexBytes, err := json.Marshal(vi)
	if err != nil {
		return nil, nil, err
	}

	encryptedIndexMapData, err := encryption.EncryptAge(indexBytes, key)
	if err != nil {
		return nil, nil, err
	}

	err = store.WriteFile(encryptedIndexMapData, indexFilePath)
	if err != nil {
		return nil, nil, err
	}

	return vi, rt, nil
}

func DeleteVault(vaultName string) error {
	if vaultName == state.State.CurrentVault {
		return fmt.Errorf("%s is current vault. Checkout to another vault first", vaultName)
	}
	vaultPath := registry.Registry.Vaults[vaultName]
	if err := store.DeleteFile(vaultPath); err != nil {
		return err
	}
	delete(registry.Registry.Vaults, vaultName)
	if err := registry.Save(); err != nil {
		return fmt.Errorf("unable to save registry: %w", err)
	}
	return nil
}
