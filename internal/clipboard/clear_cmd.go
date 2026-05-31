package clipboard

import (
	"os"

	"github.com/spf13/cobra"
)

func ClearPasswordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "clear-password",
		Hidden: true,
		Args:   cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			password := args[0]
			ClearClipboard(password)
			os.Exit(0)
		},
	}
	return cmd
}
