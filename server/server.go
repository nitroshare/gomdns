package server

import (
	"log/slog"
	"sync"
	"time"

	"github.com/nitroshare/gomdns/broadcaster"
	"github.com/nitroshare/gomdns/cache"
	"github.com/nitroshare/gomdns/dns"
	"github.com/nitroshare/gomdns/multicast"
	"github.com/nitroshare/gomdns/syncpoint"
	"github.com/nitroshare/gomdns/watcher"
)

var (
	syncAdd        = syncpoint.New()
	syncAddSuccess = syncpoint.New()
)

// Server sends and receives packets on the system's multicast interfaces.
type Server struct {

	// Received sends when a message is received.
	Received *broadcaster.Broadcaster[*dns.Message]

	logger      *slog.Logger
	cache       *cache.Cache
	watcher     *watcher.Watcher
	once        sync.Once
	chanAdded   chan multicast.Interface
	chanRemoved chan multicast.Interface
	chanRecv    chan *dns.Message
	chanSend    chan *dns.Message
	chanSendRet chan any
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
			syncAdd.Trigger()
			l, err := newServerListener(s.logger, v, s.chanRecv)
			if err != nil {
				s.logger.Warn(err.Error())
				continue
			}
			ifaceMap[v.Interface().Name] = l
			syncAddSuccess.Trigger()
		case v, ok := <-s.chanRemoved:
			if !ok {
				return
			}
			delete(ifaceMap, v.Interface().Name)
		case m := <-s.chanRecv:
			s.Received.Send(m)
		case m := <-s.chanSend:
			b, err := m.Serialize()
			if err == nil {
				for _, i := range ifaceMap {
					for _, l := range i.listeners {
						if _, err := l.Listener.Write(&multicast.Packet{
							Addr: l.Addr,
							Data: b,
						}); err != nil {
							s.logger.Error(err.Error())
						}
					}
				}
			} else {
				s.logger.Error(err.Error())
			}
			s.chanSendRet <- nil
		}
	}
}

// New creates a new Server instance.
func New(cfg *Config) *Server {
	var (
		chanAdded   = make(chan multicast.Interface)
		chanRemoved = make(chan multicast.Interface)
		s           = &Server{
			Received: broadcaster.New[*dns.Message](),
			logger:   cfg.Logger,
			cache:    cfg.Cache,
			watcher: watcher.New(
				&watcher.Config{
					Interval:    30 * time.Second,
					ChanAdded:   chanAdded,
					ChanRemoved: chanRemoved,
				},
			),
			chanAdded:   chanAdded,
			chanRemoved: chanRemoved,
			chanRecv:    make(chan *dns.Message),
			chanSend:    make(chan *dns.Message),
			chanSendRet: make(chan any),
			chanClosed:  make(chan any),
		}
	)
	if s.logger == nil {
		s.logger = slog.Default()
	}
	s.logger = s.logger.With("package", "server")
	if s.cache == nil {
		s.cache = cache.New(
			&cache.Config{
				Logger: cfg.Logger,
			},
		)
	}
	go s.run()
	return s
}

// Send transmits a packet to all of the current multicast interfaces.
func (s *Server) Send(m *dns.Message) {
	s.chanSend <- m
	<-s.chanSendRet
}

// Close shuts down the server.
func (s *Server) Close() {
	s.watcher.Close()
	<-s.chanClosed
	s.Received.Close()
}
