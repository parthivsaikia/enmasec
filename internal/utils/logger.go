package utils

import (
	"io"
	"log/slog"
)

func Logger(w io.Writer) *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(w, nil))
	return logger
}
