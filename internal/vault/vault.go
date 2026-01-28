package vault

import (
	"os"
	"path/filepath"

	"github.com/parthivsaikia/enmasec/internal/service"
	"github.com/parthivsaikia/enmasec/internal/utils"
)

type Vault struct {
	Name               string
	MasterPasswordFile string
}

func (v *Vault) InitVault() error {
	enmasecDir, err := utils.GetEnmasecDirLocation()
	if err != nil {
		return utils.ErrGetEnmasecDir(err)
	}
	vaultPath := filepath.Join(enmasecDir, v.Name)
	err = os.Mkdir(vaultPath, 0o755)
	if err != nil {
		return utils.ErrCreateDir(err, vaultPath)
	}
	return nil
}

func (v *Vault) AddService(service service.Service) error {
	enmasecDir, err := utils.GetEnmasecDirLocation()
	if err != nil {
		return utils.ErrGetEnmasecDir(err)
	}
	serviceDir := filepath.Join(enmasecDir, v.Name, service.Name)
	err = os.Mkdir(serviceDir, 0o755)
	if err != nil {
		return utils.ErrCreateDir(err, serviceDir)
	}
	return nil
}
