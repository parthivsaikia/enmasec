package store

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/parthivsaikia/enmasec/internal/encryption"
	"github.com/parthivsaikia/enmasec/internal/utils"
)

func CreateVault(vaultLocation, password string, hash []byte) error {
	if err := os.MkdirAll(vaultLocation, 0o700); err != nil {
		return fmt.Errorf("unable to create vault %w", err)
	}
	keyFile := filepath.Join(vaultLocation, "key.age")
	file, err := os.Create(keyFile)
	if err != nil {
		return fmt.Errorf("unable to create key file %w", err)
	}
	defer file.Close()

	if _, err := file.Write(append(hash, []byte("\n")...)); err != nil {
		return err
	}

	data, err := encryption.EncryptAge([]byte(KEY_FILE_TEXT), password)
	if err != nil {
		return fmt.Errorf("unable to encrypt key file: %w", err)
	}

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("unable to write to file %s", file.Name())
	}
	return nil
}

func Unlock(vaultPath, password string) error {
	var hash []byte
	var keybyte []byte
	keyfile := filepath.Join(vaultPath, "key.age")
	if !utils.CheckFileExists(keyfile) {
		fmt.Println(keyfile)
		return fmt.Errorf("key file doesn't exist")
	}

	f, err := os.Open(keyfile)
	if err != nil {
		return err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		for scanner.Text() != "age-encryption.org/v1" {
			hash = append(hash, []byte(scanner.Text())...)
		}
		keybyte = append(keybyte, scanner.Bytes()...)
	}

	key, err := encryption.DecryptAge(password, keybyte)
	if err != nil {
		return fmt.Errorf("unable to decrypt key file")
	}
	if string(key) != (KEY_FILE_TEXT) {
		return fmt.Errorf("key text doesn't match")
	}
	return nil
}
