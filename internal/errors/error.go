package errors

import (
	"fmt"
)

type Error struct {
	Code        string
	Err         error
	Description string
	Remedy      string
}

func New(code string, err error, description string, remedy string) *Error {
	return &Error{
		Code:        code,
		Err:         err,
		Description: description,
		Remedy:      remedy,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf(`
		Error: %s
		Code: %s
		Description: %s
		Remedy: %s
		`, e.Err, e.Code, e.Description, e.Remedy)
}
