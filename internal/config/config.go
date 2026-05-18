package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/parthivsaikia/enmasec/internal/models"
	"github.com/parthivsaikia/enmasec/internal/store"
	"gopkg.in/yaml.v3"
)

var Config models.Config

func checkConfigFile() string {
	files := []string{"config.yaml", "config.yml"}
	var configFile string
	for _, file := range files {
		fullpath := filepath.Join(store.GetEnmasecConfigDirLocation(), file)
		if _, err := os.Stat(fullpath); err == nil {
			configFile = fullpath
		}
	}
	return configFile
}

func Init() {
	Config.CurrentVault = ""
	Config.Vaults = map[string]string{}
}

func Load() error {
	configFile := checkConfigFile()

	configData, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}
	err = yaml.Unmarshal(configData, &Config)
	if err != nil {
		return err
	}
	for k, v := range Config.Vaults {
		if !store.CheckFileExists(v) {
			delete(Config.Vaults, k)
		}
	}
	err = Save()
	if err != nil {
		return err
	}
	return nil
}

func Save() error {
	configFile := checkConfigFile()
	if !store.CheckFileExists(store.GetEnmasecConfigDirLocation()) {
		err := os.MkdirAll(store.GetEnmasecConfigDirLocation(), 0o777)
		if err != nil {
			return fmt.Errorf("permission error: %w", err)
		}
	}
	if configFile == "" {
		configFile = filepath.Join(store.GetEnmasecConfigDirLocation(), "config.yaml")
	}
	configData, err := yaml.Marshal(Config)
	if err != nil {
		return err
	}
	err = os.WriteFile(configFile, configData, 0o666)
	if err != nil {
		return fmt.Errorf("permission error from here: %w", err)
	}
	return nil
}
