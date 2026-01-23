package vault

import (
	"encoding/json"
)

type Vault struct {
	data map[string]any
}

func (v *Vault) StorePassword() error {
	json, err := json.Marshal(v.data)
	if err != nil {
		return
	}
}
