package service

import (
	"fmt"

	"github.com/parthivsaikia/enmasec/internal/cli/components"
	"github.com/parthivsaikia/enmasec/internal/config"
	"github.com/parthivsaikia/enmasec/internal/core"
	"github.com/parthivsaikia/enmasec/internal/validation"
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
			serviceName := args[0]
			if err := validation.ValidateServiceName(serviceName); err != nil {
				return fmt.Errorf("invalid service name: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			vault, err := resolveVault(cmd)
			if err != nil {
				return err
			}

			password, err := components.PasswordPrompt(fmt.Sprintf("enter master password for vault %s", vault))
			if err != nil {
				return fmt.Errorf("unable to capture password %w", err)
			}
			key, err := core.UnlockVault(vault, password)
			if err != nil {
				return fmt.Errorf("unable to unlock vault: %w", err)
			}

			if _, _, err := core.CreateService(vault, serviceName, string(key)); err != nil {
				return fmt.Errorf("unable to create service: %w", err)
			}
			fmt.Printf("created service %s successfully.", serviceName)
			return nil
		},
	}
	addCmd.Flags().String("vault", "", "vault where service needs to be added.")
	return addCmd
}

func resolveVault(cmd *cobra.Command) (string, error) {
	vault, err := cmd.Flags().GetString("vault")
	if err != nil {
		return "", err
	}
	if vault == "" {
		vault = config.Config.CurrentVault
	}
	return vault, err
}
