package clipboard

import (
	"os"
	"os/exec"
)

func SpawnBackground(password string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command(exe, "clear-password", password)
	setSysProcAttr(cmd)
	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}
