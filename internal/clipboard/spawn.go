package clipboard

import (
	"io"
	"os"
	"os/exec"
)

func SpawnBackground(password string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command(exe, "clear-password")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	setSysProcAttr(cmd)

	if err := cmd.Start(); err != nil {
		return err
	}
	if _, err := io.WriteString(stdin, password); err != nil {
		return err
	}
	if err := stdin.Close(); err != nil {
		return err
	}

	return nil
}
