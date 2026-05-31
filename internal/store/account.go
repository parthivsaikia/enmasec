package store

import (
	"fmt"
	"os"
)

func CreateAccount(accountFilePath string, encryptedAccountData []byte) error {
	if _, err := os.Create(accountFilePath); err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}
	if err := WriteFile(encryptedAccountData, accountFilePath); err != nil {
		return fmt.Errorf("unable to write to file: %w", err)
	}
	return nil
}
