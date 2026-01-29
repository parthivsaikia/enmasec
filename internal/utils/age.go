package utils

import (
	"bytes"
	"io"
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
		return ErrEncryptFile(err, filepath)
	}

	writer.Write(data)
	writer.Close()

	return nil
}

func DecryptFromFile(masterKey string, filepath string) ([]byte, error) {
	identity, err := age.NewScryptIdentity(masterKey)
	if err != nil {
		return nil, ErrCreateAgeIdentity(err)
	}

	f, err := os.Open(filepath)
	if err != nil {
		return nil, ErrOpenFile(err, filepath)
	}
	r, err := age.Decrypt(f, identity)
	if err != nil {
		return nil, ErrDecryptFile(err, filepath)
	}

	out := &bytes.Buffer{}
	if _, err := io.Copy(out, r); err != nil {
		return nil, ErrReadEncryptionFileData(err, filepath)
	}

	return out.Bytes(), nil
}
