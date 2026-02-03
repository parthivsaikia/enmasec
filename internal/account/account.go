package account

import (
	"encoding/json"

	"github.com/parthivsaikia/enmasec/internal/utils"
)

type Account struct {
	MetaData        map[string]any
	AccountLocation string
}

func (a *Account) EncryptAccountData(masterKey string) error {
	jsonData, err := json.Marshal(a.MetaData)
	if err != nil {
		return utils.ErrJSONMarshal(err)
	}
	if err := utils.EncryptIntoFile(jsonData, masterKey, a.AccountLocation); err != nil {
		return err
	}
	return nil
}
