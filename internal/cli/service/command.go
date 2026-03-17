package service

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/parthivsaikia/enmasec/internal/config"
	"github.com/parthivsaikia/enmasec/internal/store"
	"github.com/parthivsaikia/enmasec/internal/utils"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "service",
		Short: "Manage service operations ",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}
	newCmd.AddCommand(newAddCmd())
	return newCmd
}

func newAddCmd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add [service]",
		Short: "Add a new service",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if strings.Contains(args[0], "/") {
				return fmt.Errorf("service name shouldn't contain '/'")
			}

			vault, err := cmd.Flags().GetString("vault")
			if err != nil {
				return err
			}
			if vault == "" {
				vault = config.Config.CurrentVault
			}

			servicePath := filepath.Join(config.Config.Vaults[vault], args[0])
			if utils.CheckFileExists(servicePath) {
				return fmt.Errorf("service %s already exists", args[0])
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			vault, err := cmd.Flags().GetString("vault")
			if err != nil {
				return err
			}
			if vault == "" {
				vault = config.Config.CurrentVault
			}

			vaultLocation := config.Config.Vaults[vault]
			password, err := utils.PasswordPrompt(fmt.Sprintf("enter master password for vault %s", vault))
			if err != nil {
				return fmt.Errorf("unable to capture password %w", err)
			}

			if err := store.Unlock(vaultLocation, password); err != nil {
				return fmt.Errorf("unable to unlock vault: %w", err)
			}

			serviceLocation := filepath.Join(vaultLocation, args[0])

			if err := store.CreateService(serviceLocation); err != nil {
				return fmt.Errorf("unable to create service: %w", err)
			}

			fmt.Printf("created service %s successfully.", args[0])
			return nil
		},
	}
	addCmd.Flags().String("vault", "", "vault where service needs to be added.")
	return addCmd
}
