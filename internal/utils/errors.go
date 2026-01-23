package utils

import (
	"fmt"
	"strings"
)

func NewError(err error, description []string, remedy []string) error {
	descriptionStr := strings.Join(description, "||")
	remedyStr := strings.Join(remedy, "||")
	return fmt.Errorf("Error: %v\n Description: %s\n Remedies: %s\n", err, descriptionStr, remedyStr)
}
