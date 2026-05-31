package core

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/google/uuid"
	"github.com/parthivsaikia/enmasec/internal/cli/components"
	"github.com/parthivsaikia/enmasec/internal/clipboard"
	"github.com/parthivsaikia/enmasec/internal/config"
	"github.com/parthivsaikia/enmasec/internal/encryption"
	"github.com/parthivsaikia/enmasec/internal/models"
	"github.com/parthivsaikia/enmasec/internal/store"
)

func CreateAccount(vaultName, serviceName, accountName, key string) (*models.VaultIndex, *models.RuntimeIndex, error) {
	account := models.Account{
		Username:  accountName,
		Password:  "",
		Metadata:  map[string]string{},
		UpdatedAt: time.Now().String(),
	}
	if err := components.AccountCreationREPL(&account); err != nil {
		return nil, nil, err
	}
	accountData, err := toml.Marshal(account)
	if err != nil {
		return nil, nil, err
	}
	encryptedAccountData, err := encryption.EncryptAge(accountData, key)
	if err != nil {
		return nil, nil, err
	}
	indexData, err := DecryptVaultIndex(vaultName, key)
	if err != nil {
		return nil, nil, err
	}

	vaultIndex, runtimeIndex, err := OpenVaultIndex(indexData)
	if err != nil {
		return nil, nil, err
	}

	serviceId, ok := runtimeIndex.ServiceNameToID[serviceName]
	if !ok {
		return nil, nil, fmt.Errorf("service %s doesn't exist", serviceName)
	}

	accountId := uuid.New()

	if _, ok := runtimeIndex.AccountIDToName[serviceId][accountId]; ok {
		return nil, nil, fmt.Errorf("account %s already exists", accountName)
	}
	vaultPath := config.Config.Vaults[vaultName]
	accoutFilePath := filepath.Join(vaultPath, serviceId.String(), fmt.Sprintf("%s.age", accountId.String()))
	vaultIndex.Services[serviceId].Accounts[accountId] = &models.AccountEntry{
		Name: accountName,
	}
	runtimeIndex.AccountIDToName[serviceId][accountId] = accountName
	runtimeIndex.AccountNameToID[serviceId][accountName] = accountId

	// TODO: refactor to one function
	indexBytes, err := json.Marshal(vaultIndex)
	if err != nil {
		return nil, nil, err
	}

	encryptedIndexMapData, err := encryption.EncryptAge(indexBytes, key)
	if err != nil {
		return nil, nil, err
	}
	indexFilePath := filepath.Join(vaultPath, "index.age")
	err = store.WriteFile(encryptedIndexMapData, indexFilePath)
	if err != nil {
		return nil, nil, err
	}

	if err := store.CreateAccount(accoutFilePath, encryptedAccountData); err != nil {
		return nil, nil, err
	}

	return vaultIndex, runtimeIndex, nil
}

func ListAccounts(vaultName, serviceName, key string) (*models.VaultIndex, *models.RuntimeIndex, error) {
	indexData, err := DecryptVaultIndex(vaultName, key)
	if err != nil {
		return nil, nil, err
	}

	vaultIndex, runtimeIndex, err := OpenVaultIndex(indexData)
	if err != nil {
		return nil, nil, err
	}
	if _, ok := runtimeIndex.ServiceNameToID[serviceName]; !ok {
		return nil, nil, fmt.Errorf("service %s doesn't exist", serviceName)
	}

	return vaultIndex, runtimeIndex, nil
}

func GetAccount(vaultName, serviceName, accountName, key string) (string, error) {
	indexData, err := DecryptVaultIndex(vaultName, key)
	if err != nil {
		return "", err
	}

	_, runtimeIndex, err := OpenVaultIndex(indexData)
	if err != nil {
		return "", err
	}
	serviceId, ok := runtimeIndex.ServiceNameToID[serviceName]
	if !ok {
		return "", fmt.Errorf("service doesn't exist %s", serviceName)
	}
	accountId, ok := runtimeIndex.AccountNameToID[serviceId][accountName]
	if !ok {
		return "", fmt.Errorf("account doesn't exist %s", accountName)
	}
	vaultPath := config.Config.Vaults[vaultName]
	accountPath := filepath.Join(vaultPath, serviceId.String(), accountId.String()+".age")
	accountData, err := store.ReadFile(accountPath)
	if err != nil {
		return "", err
	}
	decryptedAccountData, err := encryption.DecryptAge(key, accountData)
	if err != nil {
		return "", err
	}
	var account models.Account
	err = toml.Unmarshal(decryptedAccountData, &account)
	if err != nil {
		return "", err
	}

	if err := clipboard.SpawnBackground(account.Password); err != nil {
		return "", err
	}

	return account.Password, nil
}
func ViewAccount(vaultName, serviceName, accountName, key string) (*models.Account, error) {
	indexData, err := DecryptVaultIndex(vaultName, key)
	if err != nil {
		return nil, err
	}

	_, runtimeIndex, err := OpenVaultIndex(indexData)
	if err != nil {
		return nil, err
	}
	serviceId, ok := runtimeIndex.ServiceNameToID[serviceName]
	if !ok {
		return nil, fmt.Errorf("service doesn't exist %s", serviceName)
	}
	accountId, ok := runtimeIndex.AccountNameToID[serviceId][accountName]
	if !ok {
		return nil, fmt.Errorf("account doesn't exist %s", accountName)
	}
	vaultPath := config.Config.Vaults[vaultName]
	accountPath := filepath.Join(vaultPath, serviceId.String(), accountId.String()+".age")
	accountData, err := store.ReadFile(accountPath)
	if err != nil {
		return nil, err
	}
	decryptedAccountData, err := encryption.DecryptAge(key, accountData)
	if err != nil {
		return nil, err
	}
	var account models.Account
	err = toml.Unmarshal(decryptedAccountData, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

