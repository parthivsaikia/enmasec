package account

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	"github.com/parthivsaikia/enmasec/internal/cli/components"
	"github.com/parthivsaikia/enmasec/internal/core"
	"github.com/parthivsaikia/enmasec/internal/models"
	"github.com/parthivsaikia/enmasec/internal/state"
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
	cmd.AddCommand(newUpdateCommand())
	cmd.AddCommand(NewDeleteCommand())
	cmd.AddCommand(NewViewCommand())
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

			currentVault := state.State.CurrentVault
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

			if _, _, err := core.CreateAccount(currentVault, serviceName, accountName, string(key), &account); err != nil {
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
			currentVault := state.State.CurrentVault
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
			currentVault := state.State.CurrentVault
			password, err := components.PasswordPrompt(fmt.Sprintf("enter master password for vault %s: ", currentVault))
			if err != nil {
				return err
			}
			copy, err := cmd.Flags().GetBool("copy")
			if err != nil {
				return err
			}
			key, err := core.UnlockVault(currentVault, password)
			if err != nil {
				return fmt.Errorf("unable to unlock vault %s: %w", currentVault, err)
			}
			accountPassword, err := core.GetAccount(currentVault, serviceName, accountName, string(key), copy)
			if err != nil {
				return err
			}

			if copy {
				err := clipboard.WriteAll(accountPassword)
				if err != nil {
					return fmt.Errorf("error in copying password: %w", err)
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

func newUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [service] [account]",
		Short: "update details of an account",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			if err := validation.ValidateServiceName(serviceName); err != nil {
				return fmt.Errorf("validation error : %w", err)
			}

			accountName := args[1]
			if accountName != "" {
				if err := validation.ValidateAccountName(accountName); err != nil {
					return fmt.Errorf("validation error : %w", err)
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			accountName := args[1]
			currentVault := state.State.CurrentVault
			password, err := components.PasswordPrompt(fmt.Sprintf("enter master password for vault %s: ", currentVault))
			if err != nil {
				return err
			}
			key, err := core.UnlockVault(currentVault, password)
			if err != nil {
				return fmt.Errorf("unable to unlock vault %s: %w", currentVault, err)
			}
			account, err := core.ViewAccount(currentVault, serviceName, accountName, string(key))
			if err != nil {
				return err
			}
			if err := components.AccountUpdateREPL(account); err != nil {
				return err
			}

			if err := core.UpdateAccount(currentVault, serviceName, accountName, string(key), account); err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}

func NewDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [service] [account]",
		Short: "delete an account",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			if err := validation.ValidateServiceName(serviceName); err != nil {
				return fmt.Errorf("validation error : %w", err)
			}

			accountName := args[1]
			if accountName != "" {
				if err := validation.ValidateAccountName(accountName); err != nil {
					return fmt.Errorf("validation error : %w", err)
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			currentVault := state.State.CurrentVault
			serviceName := args[0]
			accountName := args[1]
			password, err := components.PasswordPrompt(fmt.Sprintf("enter master password for vault %s: ", currentVault))
			if err != nil {
				return err
			}
			key, err := core.UnlockVault(currentVault, password)
			if err != nil {
				return fmt.Errorf("unable to unlock vault %v: %w", currentVault, err)
			}
			if err := core.DeleteAccount(currentVault, serviceName, accountName, string(key)); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func NewViewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view [service] [account]",
		Short: "view details of an account",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			if err := validation.ValidateServiceName(serviceName); err != nil {
				return fmt.Errorf("validation error : %w", err)
			}

			accountName := args[1]
			if accountName != "" {
				if err := validation.ValidateAccountName(accountName); err != nil {
					return fmt.Errorf("validation error : %w", err)
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			currentVault := state.State.CurrentVault
			serviceName := args[0]
			accountName := args[1]
			password, err := components.PasswordPrompt(fmt.Sprintf("enter master password for vault %s: ", currentVault))
			if err != nil {
				return err
			}
			key, err := core.UnlockVault(currentVault, password)
			if err != nil {
				return fmt.Errorf("unable to unlock vault %v: %w", currentVault, err)
			}
			account, err := core.ViewAccount(currentVault, serviceName, accountName, string(key))
			if err != nil {
				return err
			}
			fmt.Println(account.Username)
			fmt.Println(account.Password)
			for k, v := range account.Metadata {
				fmt.Printf("key: %s, value: %s\n", k, v)
			}
			return nil
		},
	}
	return cmd
}
