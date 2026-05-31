package clipboard

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"
)

func WriteToClipboard(content string) error {
	if content == "" {
		return fmt.Errorf("can't write empty string to clipboard")
	}
	if err := clipboard.WriteAll(content); err != nil {
		return fmt.Errorf("can't write to clipboard: %w", err)
	}
	return nil
}

func ClearClipboard(content string) error {
	time.Sleep(30 * time.Second)
	if content == "" {
		return fmt.Errorf("can't clear empty string from clipboard")
	}
	clipboardContent, err := clipboard.ReadAll()
	if err != nil {
		return fmt.Errorf("can't read from clipboard: %w", err)
	}
	if content == clipboardContent {
		err = clipboard.WriteAll("")
		if err != nil {
			return err
		}
	}
	return nil
}
