package encryption

import (
	"github.com/parthivsaikia/enmasec/internal/errors"
)

func ErrCreateAgeRecipient(err error) error {
	return errors.New("ENC-001", err, "unable to create age recepient", "make sure that you have enough permissions.")
}
