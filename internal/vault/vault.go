package vault

import (
	"os"
	"path/filepath"

	"github.com/parthivsaikia/enmasec/internal/service"
	"github.com/parthivsaikia/enmasec/internal/utils"
)

const (
	KEY_FILE_TEXT = "Vault is now unlocked!!!!"
)

type Vault struct {
	Name     string
	Unlocked bool
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

	keyFile := filepath.Join(enmasecDir, v.Name, "key.age")
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

func (v *Vault) Unlock(masterPassword string) error {
	enmasecDir, err := utils.GetEnmasecDirLocation()
	if err != nil {
		return utils.ErrGetEnmasecDir(err)
	}
	keyFile := filepath.Join(enmasecDir, v.Name, "key.age")
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
