package cache

import (
	"log/slog"

	"github.com/nitroshare/gomdns/dns"
)

// Config provides configuration for Cache.
type Config struct {

	// ChanQuery sends records that are about to expire so that they can be
	// queried again. This can be left nil if not desired.
	ChanQuery chan<- *dns.Record

	// ChanExpired sends on records that have expired. This can be left nil if
	// not desired.
	ChanExpired chan<- *dns.Record

	// Logger can be used to capture log messages.
	Logger *slog.Logger
}
