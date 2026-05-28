package core

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/google/uuid"
	"github.com/parthivsaikia/enmasec/internal/cli/components"
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

	if _, ok := runtimeIndex.ServiceNameToID[serviceName]; !ok {
		return nil, nil, fmt.Errorf("service %s doesn't exist", serviceName)
	}

	serviceId := runtimeIndex.ServiceNameToID[serviceName]
	accountId := uuid.New()
	vaultPath := config.Config.Vaults[vaultName]
	accoutFilePath := filepath.Join(vaultPath, serviceId.String(), fmt.Sprintf("%s.age", accountId.String()))
	vaultIndex.Services[serviceId].Accounts[accountId] = &models.AccountEntry{
		Name: accountName,
	}
	runtimeIndex.AccountIDToName[serviceId][accountId] = accountName
	runtimeIndex.AccountNameToID[serviceId][accountName] = accountId

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
	serviceId := runtimeIndex.ServiceNameToID[serviceName]
	for accountName := range runtimeIndex.AccountNameToID[serviceId] {
		fmt.Println(accountName)
	}
	return vaultIndex, runtimeIndex, nil
}
