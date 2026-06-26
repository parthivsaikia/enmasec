package clipboard

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func ClearPasswordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "clear-password",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			password, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Println(err)
				writeToFile("log.txt", []byte(err.Error()))
			}
			fmt.Println(string(password))
			writeToFile("log.txt", password)
			if err := ClearClipboard(string(password)); err != nil {
				writeToFile("log.txt", []byte(err.Error()))
				return err
			}
			os.Exit(0)
			return nil
		},
	}
	return cmd
}

func writeToFile(path string, content []byte) {
	f, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	defer f.Close()
	f.Write(content)
}
