package validation

import (
	"fmt"
	"strings"
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

// func ValidateServicePath(vault, service, password string) error {
// }
