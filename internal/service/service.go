package service

import (
	"os"
	"path/filepath"

	"github.com/parthivsaikia/enmasec/internal/account"
	"github.com/parthivsaikia/enmasec/internal/utils"
)

type Service struct {
	Name            string
	ServiceLocation string
}

func (s *Service) CreateService() error {
	err := os.Mkdir(s.ServiceLocation, 0o755)
	if err != nil {
		return utils.ErrCreateDir(err, s.ServiceLocation)
	}
	return nil
}

func (s *Service) AddAccount(userName string, password string, metadata map[string]any, masterKey string) error {
	data := account.AccountData{
		Username: userName,
		Password: password,
		MetaData: metadata,
	}

	accountLocation := filepath.Join(s.ServiceLocation, userName+".age")

	account := account.Account{
		Data:            data,
		AccountLocation: accountLocation,
	}
	if err := account.EncryptAccountData(masterKey); err != nil {
		return err
	}
	return nil
}
