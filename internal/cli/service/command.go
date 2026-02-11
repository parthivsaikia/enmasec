package service

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "service",
		Short: "Manage service operations ",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	return newCmd
}
