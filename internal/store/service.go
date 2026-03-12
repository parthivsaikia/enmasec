package store

import (
	"os"
)

func CreateService(path string) error {
	if err := os.Mkdir(path, 0o700); err != nil {
		return err
	}
	return nil
}
