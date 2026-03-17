package vault

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/parthivsaikia/enmasec/internal/config"
	"github.com/parthivsaikia/enmasec/internal/encryption"
	"github.com/parthivsaikia/enmasec/internal/store"
	"github.com/parthivsaikia/enmasec/internal/utils"
	"github.com/spf13/cobra"
)

var (
	purple = lipgloss.Color("99")
	gray   = lipgloss.Color("245")
	red    = lipgloss.Red

	headerStyle  = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
	cellStyle    = lipgloss.NewStyle().Padding(0, 1)
	oddRowStyle  = cellStyle.Foreground(gray)
	evenRowStyle = cellStyle.Foreground(red)
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
			if vaultName == "" {
				return fmt.Errorf("error: vault name not provided")
			}

			if location, ok := config.Config.Vaults[vaultName]; ok {
				return fmt.Errorf("error: vault with name %s already exists in %s", vaultName, location)
			}

			if strings.ContainsAny(vaultName, "/\\") {
				return fmt.Errorf("vaultname can't contain / or \\")
			}

			dir, err := cmd.Flags().GetString("dir")
			if err != nil {
				return err
			}

			if dir == "" {
				dir = utils.GetEnmasecDirLocation()
			}

			vaultLocation := filepath.Join(dir, vaultName)

			if utils.CheckFileExists(vaultLocation) {
				return fmt.Errorf("vault %s already exists at %s", vaultName, vaultLocation)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			password, err := utils.PasswordPrompt("Set your master password: ")
			if err != nil {
				return err
			}
			if !utils.CheckPasswordValid(password) {
				return fmt.Errorf("password not strong enough")
			}
			confirmPassword, err := utils.PasswordPrompt("Enter master password again: ")
			if err != nil {
				return err
			}
			if password != confirmPassword {
				return fmt.Errorf("passwords don't match")
			}

			salt := rand.Text()

			hashedPassword := encryption.ArgonHash([]byte(password), []byte(salt))
			vaultName := args[0]
			dir, err := cmd.Flags().GetString("dir")
			if err != nil {
				return err
			}

			if dir == "" {
				dir = utils.GetEnmasecDirLocation()
			}

			vaultLocation := filepath.Join(dir, vaultName)

			err = store.CreateVault(vaultLocation, string(hashedPassword), []byte(salt))
			if err != nil {
				return err
			}

			config.Config.CurrentVault = vaultName
			config.Config.Vaults[vaultName] = vaultLocation
			if err := config.Save(); err != nil {
				return fmt.Errorf("couldn't save config: %w", err)
			}

			fmt.Printf("Created vault at %s", vaultLocation)

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
			vaultLocation := config.Config.Vaults[vaultName]

			if !utils.CheckFileExists(vaultLocation) {
				return fmt.Errorf("vault %s doesn't exist", vaultName)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultName := args[0]
			vaultPath := config.Config.Vaults[vaultName]
			password, err := utils.PasswordPrompt(fmt.Sprintf("Enter password for vault %s", vaultName))
			if err != nil {
				return err
			}
			if err := store.Unlock(vaultPath, password); err != nil {
				return fmt.Errorf("unable to unlock vault: %w", err)
			}

			if _, ok := config.Config.Vaults[vaultName]; !ok {
			} else {
				config.Config.CurrentVault = vaultName
				if err := config.Save(); err != nil {
					return err
				}
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
			var rows [][]string
			var currentVaultRow int
			for k, v := range config.Config.Vaults {
				rows = append(rows, []string{k, v})
				if k == config.Config.CurrentVault {
					currentVaultRow = len(rows) - 1
				}
			}
			t := table.New().
				Border(lipgloss.NormalBorder()).
				Headers("Name", "Location").
				StyleFunc(func(row, col int) lipgloss.Style {
					switch {
					case row == table.HeaderRow:
						return headerStyle
					case row == currentVaultRow:
						return evenRowStyle
					default:
						return oddRowStyle
					}
				}).
				Rows(rows...)
			if _, err := lipgloss.Println(t); err != nil {
				return err
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
			vaultLocation := config.Config.Vaults[vaultName]
			if !utils.CheckFileExists(vaultLocation) {
				return fmt.Errorf("vault %s doesn't exist", vaultName)
			}

			newDir, err := cmd.Flags().GetString("dir")
			if err != nil {
				return err
			}

			if newDir != "" {
				if !utils.CheckFileExists(newDir) {
					fmt.Println(newDir)
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

			if !utils.CheckPasswordValid(newPassword) {
				return fmt.Errorf("password not strong enough")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultName := args[0]
			vaultLocation := config.Config.Vaults[vaultName]
			password, err := utils.PasswordPrompt(fmt.Sprintf("Enter master password for vault %s", vaultName))
			if err != nil {
				return err
			}
			if err := store.Unlock(config.Config.Vaults[vaultName], password); err != nil {
				return fmt.Errorf("unable to unlock vault %s", vaultName)
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
				if !utils.CheckPasswordValid(newPassword) {
					return fmt.Errorf("password is not strong enough")
				}
				f := filepath.Join(newVaultLocation, "key.age")
				data, err := encryption.EncryptAge([]byte(store.KEY_FILE_TEXT), newPassword)
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
