package validation

import (
	"fmt"
	"strings"
	"unicode"

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
