package service

import (
	"os"

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
	account := account.Account{
		MetaData: metadata,
	}
	if err := account.EncryptAccountData(masterKey); err != nil {
		return err
	}
	return nil
}
