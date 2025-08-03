package server

import (
	"log/slog"
)

// Config provides configuration for Server.
type Config struct {

	// Logger can be used to capture log messages.
	Logger *slog.Logger
}
