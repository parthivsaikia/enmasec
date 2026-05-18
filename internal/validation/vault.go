package validation

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/parthivsaikia/enmasec/internal/config"
	"github.com/parthivsaikia/enmasec/internal/store"
)

func ValidateVault(vaultName, dir string) error {
	if vaultName == "" {
		return fmt.Errorf("vault name cannot be empty")
	}
	if strings.Contains(vaultName, "/\\") {
		return fmt.Errorf("vault name cannot contain / or \\")
	}
	if location, ok := config.Config.Vaults[vaultName]; ok {
		return fmt.Errorf("vault with name %s already exists in %s", vaultName, location)
	}
	vaultLocation := filepath.Join(dir, vaultName)
	if store.CheckFileExists(vaultLocation) {
		return fmt.Errorf("vault %s already exists at %s", vaultName, vaultLocation)
	}
	return nil
}

func CheckPasswordValid(password string) bool {
	if len(password) < 8 {
		return false
	}

	var hasLower, hasUpper, hasDigit, hasSpecial bool

	for _, ch := range password {
		switch {
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*", ch):
			hasSpecial = true
		}
	}

	return hasLower && hasUpper && hasDigit && hasSpecial
}
