package cmd

import (
	"fmt"
	"os"

	"github.com/parthivsaikia/enmasec/internal/cli/service"
	"github.com/parthivsaikia/enmasec/internal/cli/vault"
	"github.com/parthivsaikia/enmasec/internal/config"
	"github.com/parthivsaikia/enmasec/internal/utils"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "enmasec",
	Short: "Manage passwords from the terminal.",
	Long:  `Enmasec is a command line utility to manage passwords locally.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(vault.NewCommand())
	rootCmd.AddCommand(service.NewCommand())
	logger := utils.Logger(os.Stdout)
	config.Init()
	if err := config.Load(); err != nil {
		logger.Error(fmt.Sprintf("couldn't load config: %v", err))
	}
}
