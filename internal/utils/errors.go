package utils

import (
	"fmt"
)

func ErrCreateAgeRecipient(err error) error {
	descriptionArr := []string{
		"Couldn't create age recepient.",
		"Invalid master key.",
	}
	remedyArr := []string{
		"Make sure the master key is correct.",
	}
	return NewError(err, descriptionArr, remedyArr)
}

func ErrCreateAgeIdentity(err error) error {
	descriptionArr := []string{
		"Couldn't create age recipient",
		"Invalid master key",
	}
	remedyArr := []string{
		"Make sure the master key is correct",
	}
	return NewError(err, descriptionArr, remedyArr)
}

func ErrCreateFile(err error, filepath string) error {
	descriptionArr := []string{
		fmt.Sprintf("Couldn't create file: %s", filepath),
		"File creation failed.",
	}
	remedyArr := []string{
		"Check if the directory exists and you have write permissions.",
		"Verify the filepath is valid and doesn't contain illegal characters.",
	}
	return NewError(err, descriptionArr, remedyArr)
}

func ErrOpenFile(err error, filepath string) error {
	descriptionArr := []string{
		fmt.Sprintf("Couldn't open file %s", filepath),
		"Failed to open file",
	}
	remedyArr := []string{
		"Check if you have enough permissions to open this file.",
		"Verify the filepath is valid and doesn't contain illegal characters.",
	}
	return NewError(err, descriptionArr, remedyArr)
}

func ErrEncryptFile(err error) error {
	descriptionArr := []string{
		"Couldn't encrypt file.",
		"Encryption process failed.",
	}
	remedyArr := []string{
		"Verify the recipient/public key is valid.",
		"Ensure the file is accessible and not corrupted.",
	}
	return NewError(err, descriptionArr, remedyArr)
}

func ErrJSONMarshal(err error) error {
	return NewError(err, []string{
		"Invalid object format",
		"Unsupported data",
	}, []string{
		"Make sure to input a valid JSON object",
	})
}

func ErrGetHomeDir(err error) error {
	descriptionArr := []string{
		"Couldn't get home directory.",
		"Failed to retrieve user's home directory path.",
	}
	remedyArr := []string{
		"Verify the HOME environment variable is set.",
		"Check user account permissions and profile configuration.",
	}
	return NewError(err, descriptionArr, remedyArr)
}

func ErrGetEnmasecDir(err error) error {
	descriptionArr := []string{
		"Couldn't get enmasec directory.",
		"Failed to retrieve user's enmasec directory path.",
	}
	remedyArr := []string{
		"Check user account permissions and profile configuration.",
		"Initialise vault before adding service.",
		"Create /home/user/.enmasec directory manually.",
	}
	return NewError(err, descriptionArr, remedyArr)
}

func ErrCreateDir(err error, path string) error {
	descriptionArr := []string{
		fmt.Sprintf("Failed to create directory %s", path),
		"Directory creation failed.",
	}
	remedyArr := []string{
		"Make sure you have sufficient permissions.",
		"Verify the directory path is valid",
	}
	return NewError(err, descriptionArr, remedyArr)
}
