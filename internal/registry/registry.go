package registry

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/parthivsaikia/enmasec/internal/models"
	"github.com/parthivsaikia/enmasec/internal/store"
	"gopkg.in/yaml.v3"
)

var Registry models.Registry

func Init() {
	Registry.Vaults = make(map[string]string)
}

func Load() error {
	pruned := false
	enmasecDir := store.GetEnmasecDirLocation()
	registryFile := filepath.Join(enmasecDir, "registry.yaml")
	if _, err := os.Stat(registryFile); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	data, err := store.ReadFile(registryFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &Registry)
	if err != nil {
		return err
	}
	for k, v := range Registry.Vaults {
		if !store.CheckFileExists(v) {
			delete(Registry.Vaults, k)
			pruned = true
		}
	}
	if pruned {
		return Save()
	}
	return nil
}

func Save() error {
	enmasecDir := store.GetEnmasecDirLocation()
	registryFile := filepath.Join(enmasecDir, "registry.yaml")
	if !store.CheckFileExists(registryFile) {
		f, err := os.Create(registryFile)
		if err != nil {
			return fmt.Errorf("unable to create registry file")
		}
		defer f.Close()
	}
	registryData, err := yaml.Marshal(Registry)
	if err != nil {
		return err
	}
	err = store.WriteFile(registryData, registryFile)
	if err != nil {
		return err
	}
	return nil
}
