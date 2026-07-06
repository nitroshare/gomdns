package server

import (
	"net"
	"time"

	"github.com/nitroshare/gomulticast"
)

// Server sends and receives packets on the system's multicast interfaces.
type Server struct {
	watcher     *gomulticast.Watcher
	chanAdded   chan gomulticast.Interface
	chanRemoved chan gomulticast.Interface
	chanClosed  chan any
}

func (s *Server) run() {
	defer close(s.chanClosed)
	ifaceMap := map[string]*net.Interface{}
	for {
		select {
		case i, ok := <-s.chanAdded:
			if !ok {
				return
			}
			ifaceMap[i.Interface().Name] = i.Interface()
		case i, ok := <-s.chanRemoved:
			if !ok {
				return
			}
			delete(ifaceMap, i.Interface().Name)
		}
	}
}

// New creates a new Server instance.
func New() *Server {
	s := &Server{
		chanAdded:   make(chan gomulticast.Interface),
		chanRemoved: make(chan gomulticast.Interface),
		chanClosed:  make(chan any),
	}
	s.watcher = gomulticast.NewWatcher(
		&gomulticast.WatcherConfig{
			Interval:    2 * time.Minute,
			ChanAdded:   s.chanAdded,
			ChanRemoved: s.chanRemoved,
		},
	)
	go s.run()
	return s
}

// Send transmits a packet to all of the current multicast interfaces.
func (s *Server) Send() {
	//...
}

// Close shuts down the server.
func (s *Server) Close() {
	s.watcher.Close()
	<-s.chanClosed
}
