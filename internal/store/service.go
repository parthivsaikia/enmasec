package store

import (
	"os"

	"github.com/parthivsaikia/enmasec/internal/service"
)

func CreateService(path string) error {
	uuid := service.GenerateUUID()

	if err := os.Mkdir(path, 0o700); err != nil {
		return err
	}
	return nil
}
