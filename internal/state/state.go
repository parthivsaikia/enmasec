package state

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/parthivsaikia/enmasec/internal/models"
	"github.com/parthivsaikia/enmasec/internal/registry"
	"github.com/parthivsaikia/enmasec/internal/store"
	"gopkg.in/yaml.v3"
)

var State models.State

func Init() {
	State.CurrentVault = ""
}

func Load() error {
	stateDir := store.GetEnmasecStateDirLocation()
	stateFile := filepath.Join(stateDir, "state.yaml")
	if !store.CheckFileExists(stateFile) {
		var vaultNameArr []string
		vaults := registry.Registry.Vaults
		if len(vaults) == 0 {
			State.CurrentVault = ""
		} else {
			for v := range vaults {
				vaultNameArr = append(vaultNameArr, v)
			}
			State.CurrentVault = vaultNameArr[0]
		}

	}
	return nil
}

func Save() error {
	stateDir := store.GetEnmasecStateDirLocation()
	stateFile := filepath.Join(stateDir, "state.yaml")
	if !store.CheckFileExists(stateDir) {
		err := os.MkdirAll(stateDir, 0o777)
		if err != nil {
			return fmt.Errorf("permission error: %w", err)
		}
	}
	stateData, err := yaml.Marshal(State)
	if err != nil {
		return err
	}
	err = os.WriteFile(stateFile, stateData, 0o600)
	if err != nil {
		return err
	}
	return nil
}
