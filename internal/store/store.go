package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/parthivsaikia/enmasec/internal/encryption"
	"github.com/parthivsaikia/enmasec/internal/utils"
)

var KEY_FILE_TEXT = "This vault is locked"

// type Store struct {
// 	VaultPath string
// 	masterKey string
// }

// func New(vaultPath string) *Store {
// 	return &Store{
// 		VaultPath: vaultPath,
// 	}
// }

func Unlock(vaultPath, password string) error {
	keyfile := filepath.Join(vaultPath, "key.age")
	if !utils.CheckFileExists(keyfile) {
		fmt.Println(keyfile)
		return fmt.Errorf("key file doesn't exist")
	}
	keybyte, err := os.ReadFile(keyfile)
	if err != nil {
		return fmt.Errorf("unable to read file, %w", err)
	}
	key, err := encryption.Decrypt(password, keybyte)
	if err != nil {
		return fmt.Errorf("unable to decrypt key file")
	}
	if string(key) != (KEY_FILE_TEXT) {
		return fmt.Errorf("key text doesn't match")
	}
	return nil
}
