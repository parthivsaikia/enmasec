package vault

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/parthivsaikia/enmasec/internal/cli/components"
	"github.com/parthivsaikia/enmasec/internal/config"
	"github.com/parthivsaikia/enmasec/internal/core"
	"github.com/parthivsaikia/enmasec/internal/encryption"
	"github.com/parthivsaikia/enmasec/internal/store"
	"github.com/parthivsaikia/enmasec/internal/validation"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vault",
		Short: "Manage vault operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.Help()
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.AddCommand(newInitCommand())
	cmd.AddCommand(newCheckoutCommand())
	cmd.AddCommand(newListCommand())
	cmd.AddCommand(newUpdateCommand())
	cmd.PersistentFlags().String("dir", "", "add custom location for vault")
	return cmd
}

func newInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [vault name]",
		Short: "Initialize a new vault",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			vaultName := args[0]
			dir, err := cmd.Flags().GetString("dir")
			if err != nil {
				return err
			}
			if dir == "" {
				dir = store.GetEnmasecDirLocation()
			}

			// validate the name of the vault
			if err := validation.ValidateVaultName(vaultName); err != nil {
				return fmt.Errorf("validation error: %w", err)
			}
			vaultLocation := filepath.Join(dir, vaultName)
			if !store.CheckFileExists(vaultLocation) {
				return fmt.Errorf("vault doesn't exist")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			password, err := components.PasswordPrompt("Set your master password: ")
			if err != nil {
				return err
			}
			if !validation.CheckPasswordValid(password) {
				return fmt.Errorf("validation error: password not strong enough")
			}
			confirmPassword, err := components.PasswordPrompt("Enter master password again: ")
			if err != nil {
				return err
			}
			if password != confirmPassword {
				return fmt.Errorf("passwords don't match")
			}

			vaultName := args[0]
			dir, err := cmd.Flags().GetString("dir")
			if err != nil {
				return err
			}
			if dir == "" {
				dir = store.GetEnmasecDirLocation()
			}

			err = core.CreateVault(dir, vaultName, password)
			if err != nil {
				return fmt.Errorf("unable to create vault: %w", err)
			}

			fmt.Printf("Created vault %s at %s", vaultName, dir)
			return nil
		},
	}
	return cmd
}

func newCheckoutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checkout [vault name]",
		Short: "Checkout to another vault",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			vaultName := args[0]
			if err := validation.ValidateVaultLocationFromConfig(vaultName); err != nil {
				return fmt.Errorf("validation error: %w", err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultName := args[0]
			password, err := components.PasswordPrompt(fmt.Sprintf("Enter password for vault %s: ", vaultName))
			if err != nil {
				return err
			}
			if _, err := core.UnlockVault(vaultName, password); err != nil {
				return fmt.Errorf("unable to unlock vault: %w", err)
			}

			if err := core.CheckoutVault(vaultName); err != nil {
				return fmt.Errorf("unable to checkout vault: %w", err)
			}
			fmt.Printf("Switched to vault %s", vaultName)
			return nil
		},
	}
	return cmd
}

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all the available vaults",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := core.ListVaults(); err != nil {
				return fmt.Errorf("unable to list vaults: %w", err)
			}
			return nil
		},
	}
	return cmd
}

func newUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [vault name]",
		Short: "Update name, password or location of a vault",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			vaultName := args[0]
			if err := validation.ValidateVaultLocationFromConfig(vaultName); err != nil {
				return err
			}

			newDir, err := cmd.Flags().GetString("dir")
			if err != nil {
				return err
			}

			if newDir != "" {
				if !store.CheckFileExists(newDir) {
					return fmt.Errorf("directory %s doesn't exist", newDir)
				}
			}

			newName, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			if _, ok := config.Config.Vaults[newName]; ok {
				return fmt.Errorf("vault with name %s already exist", newName)
			}

			newPassword, err := cmd.Flags().GetString("password")
			if err != nil {
				return err
			}

			if !validation.CheckPasswordValid(newPassword) {
				return fmt.Errorf("password not strong enough")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultName := args[0]
			vaultLocation := config.Config.Vaults[vaultName]
			password, err := components.PasswordPrompt(fmt.Sprintf("Enter master password for vault %s", vaultName))
			if err != nil {
				return err
			}
			key, err := core.UnlockVault(vaultName, password)
			if err != nil {
				return err
			}
			newDir, err := cmd.Flags().GetString("dir")
			if err != nil {
				return err
			}

			if newDir == "" {
				newDir = filepath.Dir(vaultLocation)
			}

			newName, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			if newName == "" {
				newName = vaultName
			}

			newVaultLocation := filepath.Join(newDir, newName)

			if newVaultLocation != "" {
				if err := os.Rename(vaultLocation, newVaultLocation); err != nil {
					return err
				}
			}

			newPassword, err := cmd.Flags().GetString("password")
			if err != nil {
				return err
			}

			if newPassword != "" {
				if !validation.CheckPasswordValid(newPassword) {
					return fmt.Errorf("password is not strong enough")
				}
				f := filepath.Join(newVaultLocation, "key.age")
				data, err := encryption.EncryptAge(key, newPassword)
				if err != nil {
					return err
				}
				err = os.WriteFile(f, data, 0o666)
				if err != nil {
					return err
				}
			}

			config.Config.Vaults[newName] = newVaultLocation
			if err := config.Save(); err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().String("password", "", "change password of the vault.")
	cmd.Flags().String("name", "", "change name of the vault.")
	return cmd
}
