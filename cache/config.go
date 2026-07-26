package cache

import (
	"log/slog"
)

// Config provides configuration for Cache.
type Config struct {

	// Logger can be used to capture log messages.
	Logger *slog.Logger
}
