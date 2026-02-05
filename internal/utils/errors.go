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

func ErrWriteFile(err error, filepath string) error {
	descriptionArr := []string{
		fmt.Sprintf("Error in writing to file %s", filepath),
	}
	remedyArr := []string{
		"Make sure you have write permissions.",
		"Make sure the file name is correct.",
	}
	return NewError(err, descriptionArr, remedyArr)
}

func ErrEncryptFile(err error, filepath string) error {
	descriptionArr := []string{
		fmt.Sprintf("Couldn't encrypt file %s.", filepath),
		"Encryption process failed.",
	}
	remedyArr := []string{
		"Verify the recipient/public key is valid.",
		"Ensure the file is accessible and not corrupted.",
	}
	return NewError(err, descriptionArr, remedyArr)
}

func ErrDecryptFile(err error, filepath string) error {
	descriptionArr := []string{
		fmt.Sprintf("Couldn't decrypt file %s.", filepath),
		"Decryption process failed.",
	}
	remedyArr := []string{
		"Verify the recipient/public key is valid.",
		"Ensure the file is accessible and not corrupted.",
	}
	return NewError(err, descriptionArr, remedyArr)
}

func ErrReadEncryptionFileData(err error, filepath string) error {
	descriptionArr := []string{
		fmt.Sprintf("File %s is decrypted but can't read its content", filepath),
		"Error in reading encrypted data",
	}
	remedyArr := []string{
		"Check if the file has proper read permissions.",
		"Verify the file isn't corrupted or locked by another process.",
		"Ensure sufficient memory is available to read the file.",
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

func ErrUnlockVault(vaultName string) error {
	err := (fmt.Errorf("failed to unlock vault %s", vaultName))
	descriptionArr := []string{
		"Failed to open vault",
		"Decrypted text doesn't match with key.",
	}
	remedyArr := []string{
		"Make sure that you have enough permissions to read the key file",
	}
	return NewError(err, descriptionArr, remedyArr)
}

func ErrRenameDir(err error, oldPath, newPath string) error {
	descriptionArr := []string{
		fmt.Sprintf("Cannot rename %s to %s", oldPath, newPath),
		"Cannot perform renaming operation",
	}
	remedyArr := []string{
		"Make sure you have enough permissions to write",
		fmt.Sprintf("Make sure the filepath %s exists", oldPath),
	}
	return NewError(err, descriptionArr, remedyArr)
}
