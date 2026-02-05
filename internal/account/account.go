package account

import (
	"encoding/json"
	"fmt"

	"github.com/parthivsaikia/enmasec/internal/utils"
)

type AccountData struct {
	Username string         `json:"username"`
	Password string         `json:"password"`
	MetaData map[string]any `json:"metadata"`
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

func (a *Account) ReadAccountData(masterKey string) (*AccountData, error) {
	data, err := utils.DecryptFromFile(masterKey, a.AccountLocation)
	if err != nil {
		return nil, err
	}
	accoundData := AccountData{}
	err = json.Unmarshal(data, &accoundData)
	if err != nil {
		return nil, fmt.Errorf("error in unmarshalling json to struct")
	}
	return &accoundData, err
}
