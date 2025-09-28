package logger

import (
	"log/slog"
	"os"
)

// services logger
var Logger *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
