package browser

import (
	"log/slog"

	"github.com/nitroshare/gomdns/cache"
	"github.com/nitroshare/gomdns/server"
)

// Config provides configuration data for Browser.
type Config struct {

	// Type indicates the service type you are browsing.
	Type string

	// Cache provides a cache to use for managing records. This can be left
	// nil if not needed but is highly recommended if you are using multiple
	// browsers or using a provider.
	Cache *cache.Cache

	// Server provides a server to use for sending queries and listening for
	// responses.
	Server *server.Server

	// Logger can be used to capture log messages.
	Logger *slog.Logger
}
