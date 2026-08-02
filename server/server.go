package server

import (
	"errors"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/nitroshare/gomdns"
	"github.com/nitroshare/gomdns/broadcaster"
	"github.com/nitroshare/gomdns/cache"
	"github.com/nitroshare/gomdns/dns"
	"github.com/nitroshare/gomdns/syncpoint"
	"github.com/nitroshare/gomdns/vnet"
	"github.com/nitroshare/gomdns/watcher"
)

var (
	syncAdd           = syncpoint.New()
	syncAddSuccess    = syncpoint.New()
	syncRemoveSuccess = syncpoint.New()
)

// Server sends and receives packets on the system's multicast interfaces.
type Server struct {

	// Received sends when a message is received.
	Received *broadcaster.Broadcaster[*dns.Message]

	logger      *slog.Logger
	cache       *cache.Cache
	watcher     *watcher.Watcher
	once        sync.Once
	chanAdded   chan vnet.Interface
	chanRemoved chan vnet.Interface
	chanRecv    chan *dns.Message
	chanSend    chan *dns.Message
	chanSendRet chan any
	chanClosed  chan any
}

func (s *Server) createConns(i vnet.Interface) ([]*serverConn, error) {
	var (
		listConn = []*serverConn{}
		flags    = i.Flags()
	)
	if flags&net.FlagUp == 0 ||
		flags&net.FlagRunning == 0 ||
		flags&net.FlagMulticast == 0 {
		return nil, errors.New("interface does not support multicast")
	}
	v, err := i.Addrs()
	if err != nil {
		return nil, err
	}
	for _, a := range v {
		ipnet, ok := a.(*net.IPNet)
		if !ok {
			continue
		}
		ip := ipnet.IP
		switch {
		case ip.To4() != nil:
			if ip.IsLoopback() {
				continue
			}
			c, err := newServerConn(i, "udp4", gomdns.IPv4Address, s.chanRecv)
			if err != nil {
				s.logger.Warn(err.Error())
				continue
			}
			listConn = append(listConn, c)
		default:
			if ip.IsLinkLocalUnicast() {
				continue
			}
			c, err := newServerConn(i, "udp4", gomdns.IPv6Address, s.chanRecv)
			if err != nil {
				s.logger.Warn(err.Error())
				continue
			}
			listConn = append(listConn, c)
		}
	}
	if len(listConn) == 0 {
		return nil, errors.New("no supported addresses")
	}
	return listConn, nil
}

func (s *Server) run() {
	defer close(s.chanClosed)
	ifaceMap := map[string][]*serverConn{}
	defer func() {
		for _, a := range ifaceMap {
			for _, c := range a {
				c.Close()
			}
		}
	}()
	for {
		select {
		case v, ok := <-s.chanAdded:
			if !ok {
				return
			}
			syncAdd.Trigger()
			a, err := s.createConns(v)
			if err != nil {
				s.logger.Warn(err.Error())
				continue
			}
			ifaceMap[v.Name()] = a
			syncAddSuccess.Trigger()
		case v, ok := <-s.chanRemoved:
			if !ok {
				return
			}
			delete(ifaceMap, v.Name())
			syncRemoveSuccess.Trigger()
		case m := <-s.chanRecv:
			s.Received.Send(m)
		case m := <-s.chanSend:
			b, err := m.Serialize()
			if err == nil {
				for _, a := range ifaceMap {
					for _, c := range a {
						c.conn.WriteTo(b, c.addr)
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
		chanAdded   = make(chan vnet.Interface)
		chanRemoved = make(chan vnet.Interface)
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
