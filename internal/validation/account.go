package validation

import (
	"fmt"
	"strings"
)

func ValidateAccountName(name string) error {
	if strings.Contains(name, "/\\") {
		return fmt.Errorf("account name cannot contain / or \\")
	}
	return nil
}

func ValidateAccountMetaDataKey(key string) error {
	if key == "" {
		return fmt.Errorf("key or value can't be empty")
	}
	return nil
}
