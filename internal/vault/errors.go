package vault

import (
	"github.com/parthivsaikia/enmasec/internal/utils"
)

func ErrJsonMarshal(err error) error {
	return utils.NewError(err, []string{
		"Invalid object format",
		"Unsupported data",
	}, []string{
		"Make sure to input a valid JSON object",
	})
}
