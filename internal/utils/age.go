package utils

import (
	"os"

	"filippo.io/age"
)

func EncryptIntoFile(data []byte, masterKey string, filepath string) error {
	recipient, err := age.NewScryptRecipient(masterKey)
	if err != nil {
		return ErrCreateAgeRecipient(err)
	}

	f, err := os.Create(filepath)
	if err != nil {
		return ErrCreateFile(err, filepath)
	}

	defer f.Close()

	writer, err := age.Encrypt(f, recipient)
	if err != nil {
		return ErrEncryptFile(err)
	}

	writer.Write(data)
	writer.Close()

	return nil
}
