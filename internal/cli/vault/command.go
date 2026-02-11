package vault

import (
	"fmt"
	"path/filepath"

	"github.com/parthivsaikia/enmasec/internal/utils"
	"github.com/parthivsaikia/enmasec/internal/vault"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vault",
		Short: "Manage vault operations",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(newInitCommand())
	return cmd
}

func newInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [vault name]",
		Short: "Initialize a new vault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			password, err := utils.PasswordPrompt("Set your master password: ")
			if err != nil {
				return err
			}
			confirmPassword, err := utils.PasswordPrompt("Enter master password again: ")
			if err != nil {
				return err
			}
			if password != confirmPassword {
				return fmt.Errorf("passwords don't match")
			}

			enmasecDir, err := utils.GetEnmasecDirLocation()
			if err != nil {
				return err
			}

			newVault := vault.Vault{
				Name:          args[0],
				VaultLocation: filepath.Join(enmasecDir, args[0]),
				Unlocked:      false,
			}
			if err := newVault.InitVault(password); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
