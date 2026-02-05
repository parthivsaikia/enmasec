package vault

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/parthivsaikia/enmasec/internal/service"
	"github.com/parthivsaikia/enmasec/internal/utils"
)

const (
	KEY_FILE_TEXT = "Vault is now unlocked!!!!"
)

type Vault struct {
	Name          string
	VaultLocation string
	Unlocked      bool
}

func (v *Vault) InitVault() error {
	err := os.Mkdir(v.VaultLocation, 0o755)
	if err != nil {
		return utils.ErrCreateDir(err, v.VaultLocation)
	}

	keyFile := filepath.Join(v.VaultLocation, "key.age")
	f, err := os.Create(keyFile)
	if err != nil {
		return utils.ErrCreateFile(err, keyFile)
	}
	defer f.Close()
	_, err = f.Write([]byte(KEY_FILE_TEXT))
	if err != nil {
		return utils.ErrWriteFile(err, keyFile)
	}
	return nil
}

func (v *Vault) AddService(serviceName string) error {
	serviceDir := filepath.Join(v.VaultLocation, serviceName)

	service := service.Service{
		Name:            serviceName,
		ServiceLocation: serviceDir,
	}
	if err := service.CreateService(); err != nil {
		return err
	}
	return nil
}

func (v *Vault) Unlock(masterPassword string) error {
	keyFile := filepath.Join(v.VaultLocation, "key.age")
	key, err := utils.DecryptFromFile(masterPassword, keyFile)
	if err != nil {
		return err
	}
	if string(key) == (KEY_FILE_TEXT) {
		v.Unlocked = true
		return nil
	}
	return utils.ErrUnlockVault(v.Name)
}

func (v *Vault) UpdateService(oldService service.Service, newServiceName string) error {
	// oldServicePath = ~/enmasec/vaultName/service
	newServiceName = filepath.Join(v.VaultLocation, newServiceName)
	if err := os.Rename(oldService.Name, newServiceName); err != nil {
		return utils.ErrRenameDir(err, oldService.Name, newServiceName)
	}
	return nil
}

func (v *Vault) GetServices() ([]os.DirEntry, error) {
	entities, err := os.ReadDir(v.VaultLocation)
	serviceDirs := []os.DirEntry{}
	if err != nil {
		// TODO: replace fmt.Errorf with enmasec Errors
		return nil, fmt.Errorf("error reading dir %s", v.VaultLocation)
	}
	for _, entity := range entities {
		if !entity.IsDir() {
			continue
		}
		serviceDirs = append(serviceDirs, entity)
	}
	return serviceDirs, nil
}
