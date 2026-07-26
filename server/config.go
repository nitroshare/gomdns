package server

import (
	"log/slog"

	"github.com/nitroshare/gomdns/cache"
)

// Config provides configuration for Server.
type Config struct {

	// Logger can be used to capture log messages.
	Logger *slog.Logger

	// Cache provides a cache for DNS records.
	Cache *cache.Cache
}
