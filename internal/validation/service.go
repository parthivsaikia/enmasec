package validation

import (
	"fmt"
	"strings"
)

func ValidateServiceName(name string) error {
	if name == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	if strings.ContainsAny(name, "/\\") {
		return fmt.Errorf("service name cannot contain / or \\")
	}
	if strings.TrimSpace(name) != name {
		return fmt.Errorf("service name cannot have leading or trailing spaces")
	}
	if len(name) > 255 {
		return fmt.Errorf("service name cannot exceed 255 characters")
	}
	return nil
}
