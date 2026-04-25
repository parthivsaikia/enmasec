package validation

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/parthivsaikia/enmasec/internal/models"
	"github.com/parthivsaikia/enmasec/internal/store"
)

func ValidateServiceName(name string) error {
	if name == "" {
		return fmt.Errorf("vault or service name cannot be empty")
	}
	if strings.Contains(name, "/\\") {
		return fmt.Errorf("service name cannot contain / or \\")
	}
	return nil
}

func ValidateServicePath(vault, service, password string) error {
	servicePath := filepath.Join(vault, service)
	// decrypt the bidirectional map and check if the service already exists
	mapData, err := store.DecryptIndexFile(vault, password)
	if err != nil {
		return err
	}
	biMap, err := json.Unmarshal(mapData, models.BiMap)
	if err !
}
