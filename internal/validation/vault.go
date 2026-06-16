package validation

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/parthivsaikia/enmasec/internal/config"
	"github.com/parthivsaikia/enmasec/internal/store"
)

func ValidateVaultName(vaultName string) error {
	if vaultName == "" {
		return fmt.Errorf("vault name cannot be empty")
	}
	if strings.Contains(vaultName, "/\\") {
		return fmt.Errorf("vault name cannot contain / or \\")
	}
	return nil
}

func ValidateVaultLocationFromConfig(vaultName string) bool {
	vaultLocation, ok := config.Config.Vaults[vaultName]
	if !ok {
		return false
	}
	return store.CheckFileExists(vaultLocation)
}

func IsVaultCurrentVault(vaultName string) bool {
	return vaultName == config.Config.CurrentVault
}

func CheckPasswordValid(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be longer than or equal to 8 characters")
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

	if !(hasLower && hasUpper && hasDigit && hasSpecial) {
		return fmt.Errorf("password not strong enough")
	}
	return nil
}
