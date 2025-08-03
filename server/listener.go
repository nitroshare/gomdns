package server

import (
	"log/slog"
)

// Server sends and receives DNS messages on all available multicast
// interfaces.
type Server struct {
	logger    *slog.Logger
	chanClose chan any
}

func (s *Server) run() {
	defer close(s.chanClose)
	for {
		select {
		case <-s.chanClose:
			return
		}
	}
}

// New creates a new Server.
func New(cfg *Config) *Server {
	s := &Server{
		logger: cfg.Logger,
	}
	if s.logger == nil {
		s.logger = slog.Default()
	}
	s.logger = s.logger.With("package", "server")
	go s.run()
	return s
}

// Close shuts down the server.
func (s *Server) Close() {
	s.chanClose <- nil
	<-s.chanClose
}
