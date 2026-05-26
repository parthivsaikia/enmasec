package core

import (
	"path/filepath"

	"github.com/google/uuid"
	"github.com/parthivsaikia/enmasec/internal/models"
)

func CreateAccount(vaultName, serviceName, accountName, key string) (*models.VaultIndex, *models.RuntimeIndex, error) {
	indexData, err := DecryptVaultIndex(vaultName, key)
	if err != nil {
		return nil, nil, err
	}

	vaultIndex, runtimeIndex, err := OpenVaultIndex(indexData)
	if err != nil {
		return nil, nil, err
	}

	serviceId := runtimeIndex.ServiceNameToID[serviceName]
	accountId := uuid.New()
	accoutFilePath := filepath.Join(vaultName, serviceId.String(), accountId.String())
	vaultIndex.Services[serviceId].Accounts[accountId] = &models.AccountEntry{
		Name: accountName,
	}
	runtimeIndex.AccountIDToName[serviceId][accountId] = accountName
	runtimeIndex.AccountNameToID[serviceId][accountName] = accountId
	return vaultIndex, runtimeIndex, nil
}
