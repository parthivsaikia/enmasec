package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/parthivsaikia/enmasec/internal/cli/components"
	"github.com/parthivsaikia/enmasec/internal/config"
	"github.com/parthivsaikia/enmasec/internal/encryption"
	"github.com/parthivsaikia/enmasec/internal/models"
	"github.com/parthivsaikia/enmasec/internal/store"
	"github.com/parthivsaikia/enmasec/internal/validation"
)

func CreateVault(dir, vaultName, password string) error {
	vaultPath := filepath.Join(dir, vaultName)
	secretKey := encryption.RandomByte(32)
	encryptedKey, err := encryption.EncryptAge(secretKey, password)
	if err != nil {
		return fmt.Errorf("unable to encrypt: %w", err)
	}
	emptyIndex := &models.VaultIndex{
		Version:  1,
		Services: make(map[uuid.UUID]*models.ServiceEntry),
	}
	indexByte, err := json.Marshal(emptyIndex)
	if err != nil {
		return fmt.Errorf("unable to encrypt index data: %w", err)
	}
	encryptedIndexByte, err := encryption.EncryptAge(indexByte, string(secretKey))
	err = store.CreateVaultStore(vaultPath, password, encryptedKey, encryptedIndexByte)
	if err != nil {
		return err
	}
	config.Config.CurrentVault = vaultName
	config.Config.Vaults[vaultName] = vaultPath
	if err := config.Save(); err != nil {
		return fmt.Errorf("couldn't save config: %w", err)
	}
	return nil
}

func UnlockVault(vaultName, password string) ([]byte, error) {
	vaultPath := config.Config.Vaults[vaultName]
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
	if _, ok := config.Config.Vaults[vaultName]; !ok {
		return fmt.Errorf("vault doesn't exist")
	} else {
		config.Config.CurrentVault = vaultName
		if err := config.Save(); err != nil {
			return fmt.Errorf("unable to save config: %w", err)
		}
	}
	return nil
}

func ListVaults() error {
	var rows [][]string
	var currentVaultRow int
	for k, v := range config.Config.Vaults {
		rows = append(rows, []string{k, v})
		if k == config.Config.CurrentVault {
			currentVaultRow = len(rows) - 1
		}
	}
	if err := components.VaultTable(currentVaultRow, rows); err != nil {
		return err
	}
	return nil
}

func UpdateVault(vaultName, newVaultName, newDir, newPassword string, key []byte) error {
	vaultPath := config.Config.Vaults[vaultName]
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
		if !validation.CheckPasswordValid(newPassword) {
			return fmt.Errorf("password is not strong enough")
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

	config.Config.Vaults[newVaultName] = newVaultLocation
	config.Config.CurrentVault = newVaultName
	if newVaultName != vaultName {
		delete(config.Config.Vaults, vaultName)
	}
	if err := config.Save(); err != nil {
		return err
	}

	return nil
}

func DecryptVaultIndex(vaultName, key string) ([]byte, error) {
	vaultPath := config.Config.Vaults[vaultName]
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
	vaultPath := config.Config.Vaults[vaultName]
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
