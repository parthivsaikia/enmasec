func WriteToClipboard(content string) error {
	if content == "" {
		return fmt.Errorf("can't write empty string to clipboard")
	}
	if err := clipboard.WriteAll(content); err != nil {
		return fmt.Errorf("can't write to clipboard: %w", err)
	}
	return nil
}
