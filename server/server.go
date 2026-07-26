package server

import (
	"log/slog"
	"time"

	"github.com/nitroshare/gomdns/broadcaster"
	"github.com/nitroshare/gomdns/cache"
	"github.com/nitroshare/gomdns/dns"
	"github.com/nitroshare/gomdns/multicast"
	"github.com/nitroshare/gomdns/watcher"
)

// Server sends and receives packets on the system's multicast interfaces.
type Server struct {

	// Received sends when an incoming
	Received *broadcaster.Broadcaster[*dns.Record]

	logger      *slog.Logger
	cache       *cache.Cache
	watcher     *watcher.Watcher
	chanAdded   chan multicast.Interface
	chanRemoved chan multicast.Interface
	chanSend    chan *dns.Message
	chanClosed  chan any
}

func (s *Server) run() {
	defer close(s.chanClosed)
	ifaceMap := map[string]*serverListener{}
	defer func() {
		for _, i := range ifaceMap {
			i.Close()
		}
	}()
	for {
		select {
		case v, ok := <-s.chanAdded:
			if !ok {
				return
			}
			l, err := newServerListener(s.logger, v, s.chanSend)
			if err != nil {
				s.logger.Warn(err.Error())
				continue
			}
			ifaceMap[v.Interface().Name] = l
		case v, ok := <-s.chanRemoved:
			if !ok {
				return
			}
			delete(ifaceMap, v.Interface().Name)
		case v := <-s.chanSend:
			for _, r := range v.Records {
				s.Received.Send(r)
			}
		}
	}
}

// New creates a new Server instance.
func New(cfg *Config) *Server {
	var (
		chanAdded   = make(chan multicast.Interface)
		chanRemoved = make(chan multicast.Interface)
		s           = &Server{
			Received: broadcaster.New[*dns.Record](),
			logger:   cfg.Logger,
			cache: cache.New(
				&cache.Config{
					Logger: cfg.Logger,
				},
			),
			watcher: watcher.New(
				&watcher.Config{
					Interval:    30 * time.Second,
					ChanAdded:   chanAdded,
					ChanRemoved: chanRemoved,
				},
			),
			chanAdded:   chanAdded,
			chanRemoved: chanRemoved,
			chanSend:    make(chan *dns.Message),
			chanClosed:  make(chan any),
		}
	)
	if s.logger == nil {
		s.logger = slog.Default()
	}
	s.logger = s.logger.With("package", "server")
	go s.run()
	return s
}

// Send transmits a packet to all of the current multicast interfaces.
func (s *Server) SendAll() {
	//...
}

// Close shuts down the server.
func (s *Server) Close() {
	s.watcher.Close()
	<-s.chanClosed
}
