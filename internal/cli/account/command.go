package account

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/parthivsaikia/enmasec/internal/cli/components"
	"github.com/parthivsaikia/enmasec/internal/config"
	"github.com/parthivsaikia/enmasec/internal/core"
	"github.com/parthivsaikia/enmasec/internal/validation"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "Manage account operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.Help()
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.AddCommand(newAddCommand())
	cmd.AddCommand(newListCommand())
	cmd.AddCommand(newGetCommand())
	return cmd
}

func newAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [service] [account]",
		Short: "add new account",
		Args:  cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			if err := validation.ValidateServiceName(serviceName); err != nil {
				return fmt.Errorf("validation error : %w", err)
			}

			accountName := args[1]
			if err := validation.ValidateAccountName(accountName); err != nil {
				return fmt.Errorf("validation error : %w", err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			accountName := args[1]

			currentVault := config.Config.CurrentVault
			password, err := components.PasswordPrompt(fmt.Sprintf("enter master password for vault %s", currentVault))
			if err != nil {
				return err
			}
			key, err := core.UnlockVault(currentVault, password)
			if err != nil {
				return fmt.Errorf("unable to unlock vault %s: %w", currentVault, err)
			}
			account := models.Account{
				Username:  accountName,
				Password:  "",
				Metadata:  map[string]string{},
				UpdatedAt: time.Now().String(),
			}
			if err := components.AccountCreationREPL(&account); err != nil {
				return err
			}

			if _, _, err := core.CreateAccount(currentVault, serviceName, accountName, string(key)); err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [service]",
		Short: "list all accounts of a service",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			if err := validation.ValidateServiceName(serviceName); err != nil {
				return fmt.Errorf("validation error : %w", err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			currentVault := config.Config.CurrentVault
			password, err := components.PasswordPrompt(fmt.Sprintf("enter master password for vault %s: ", currentVault))
			if err != nil {
				return err
			}
			key, err := core.UnlockVault(currentVault, password)
			if err != nil {
				return fmt.Errorf("unable to unlock vault %s: %w", currentVault, err)
			}

			_, rt, err := core.ListAccounts(currentVault, serviceName, string(key))

			serviceId := rt.ServiceNameToID[serviceName]
			for accountName := range rt.AccountNameToID[serviceId] {
				fmt.Println(accountName)
			}
			if err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func newGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [service] [account]",
		Short: "get the password ",
		Args:  cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			if err := validation.ValidateServiceName(serviceName); err != nil {
				return fmt.Errorf("validation error : %w", err)
			}

			accountName := args[1]
			if err := validation.ValidateAccountName(accountName); err != nil {
				return fmt.Errorf("validation error : %w", err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			accountName := args[1]
			currentVault := config.Config.CurrentVault
			password, err := components.PasswordPrompt(fmt.Sprintf("enter master password for vault %s: ", currentVault))
			if err != nil {
				return err
			}
			key, err := core.UnlockVault(currentVault, password)
			if err != nil {
				return fmt.Errorf("unable to unlock vault %s: %w", currentVault, err)
			}
			accountPassword, err := core.GetAccount(currentVault, serviceName, accountName, string(key))
			if err != nil {
				return err
			}
			copy, err := cmd.Flags().GetBool("copy")
			if err != nil {
				return err
			}
			if copy {
				err := clipboard.WriteAll(accountPassword)
				if err != nil {
					return err
				}
			} else {
				fmt.Println(accountPassword)
			}
			return nil
		},
	}
	cmd.Flags().Bool("copy", false, "copy password to clipboard")
	return cmd
}
