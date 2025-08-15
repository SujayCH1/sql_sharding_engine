package services

import (
	"log/slog"
	"os"
	"time"
)

// Query struct for entire applications
type Query struct {
	QueryString string
	Timestamp   time.Time
	queryID     string
}

// temp pk of database
const KeyColumn string = "pk"

// services logger
var Logger *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
