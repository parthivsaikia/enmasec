package validation

import (
	"fmt"
	"strings"

	"github.com/parthivsaikia/enmasec/internal/config"
)

func ValidateVaultName(vaultName string) error {
	if vaultName == "" {
		return fmt.Errorf("vault name cannot be empty")
	}
	if strings.Contains(vaultName, "/\\") {
		return fmt.Errorf("vault name cannot contain / or \\")
	}
	if location, ok := config.Config.Vaults[vaultName]; ok {
		return fmt.Errorf("vault with name %s already exists in %s", vaultName, location)
	}
	return nil
}
