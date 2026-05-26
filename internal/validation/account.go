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
