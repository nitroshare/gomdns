package browser

import (
	"log/slog"

	"github.com/nitroshare/gomdns/cache"
	"github.com/nitroshare/gomdns/server"
)

// Browser discovers services running on the local network.
type Browser struct {
	logger      *slog.Logger
	cache       *cache.Cache
	cacheClose  bool
	server      *server.Server
	serverClose bool
	chanClose   chan any
	chanClosed  chan any
}

func (b *Browser) run() {
	defer close(b.chanClosed)
	for {
		select {
		case <-b.chanClose:
			return
		}
	}
}

// New creates a new Browser instance.
func New(cfg *Config) *Browser {
	b := &Browser{
		logger:     cfg.Logger,
		cache:      cfg.Cache,
		server:     cfg.Server,
		chanClose:  make(chan any),
		chanClosed: make(chan any),
	}
	if b.cache == nil {
		b.cache = cache.New(
			&cache.Config{
				Logger: b.logger,
			},
		)
		b.cacheClose = true
	}
	if b.server == nil {
		b.server = server.New()
		b.serverClose = true
	}
	if b.logger == nil {
		b.logger = slog.Default()
	}
	b.logger = b.logger.With("package", "cache")
	go b.run()
	return b
}

// Close shuts down the browser.
func (b *Browser) Close() {
	close(b.chanClose)
	<-b.chanClosed
	if b.serverClose {
		b.server.Close()
	}
	if b.cacheClose {
		b.cache.Close()
	}
}
