package messages

import "github.com/parthivsaikia/enmasec/internal/models"

type VaultCreateMsg *models.Vault

type VaultAddMsg []*models.Vault

type ErrMsg error
