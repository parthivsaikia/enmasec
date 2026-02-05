package account

import (
	"encoding/json"

	"github.com/parthivsaikia/enmasec/internal/utils"
)

type AccountData struct {
	Username string
	Password string
	MetaData map[string]any
}

type Account struct {
	Data            AccountData
	AccountLocation string
}

func (a *Account) EncryptAccountData(masterKey string) error {
	jsonData, err := json.Marshal(a.Data)
	if err != nil {
		return utils.ErrJSONMarshal(err)
	}
	if err := utils.EncryptIntoFile(jsonData, masterKey, a.AccountLocation); err != nil {
		return err
	}
	return nil
}
