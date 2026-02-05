package service

import (
	"fmt"
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

func (s *Service) GetAccounts() ([]os.DirEntry, error) {
	entities, err := os.ReadDir(s.ServiceLocation)
	if err != nil {
		// TODO: replace fmt.Errorf with enmasec errors
		return nil, fmt.Errorf("error in reading directory %v", err)
	}
	accounts := []os.DirEntry{}
	for _, entity := range entities {
		if entity.IsDir() {
			continue
		} else {
			accounts = append(accounts, entity)
		}
	}
	return accounts, nil
}

func (s *Service) UpdateAccountName(oldAccount account.Account, newName string) error {
	newAccountLocation := filepath.Join(s.ServiceLocation, newName)
	if err := os.Rename(oldAccount.AccountLocation, newAccountLocation); err != nil {
		return fmt.Errorf("Couldn't rename file %s %v", oldAccount.AccountLocation, err)
	}
	return nil
}
