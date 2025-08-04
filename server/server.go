package server

import (
	"context"
	"log/slog"
	"net"
	"sync"
	"syscall"
	"time"

	"github.com/nitroshare/mocktime"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

const mDNSListenAddr = ":5353"

var (
	mDNSMulticastIPv4Addr = &net.UDPAddr{
		IP: net.ParseIP("224.0.0.251"),
	}
	mDNSMulticastIPv6Addr = &net.UDPAddr{
		IP: net.ParseIP("ff02::fb"),
	}
	listenConfig = net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				syscall.SetsockoptInt(
					syscall.Handle(fd),
					syscall.SOL_SOCKET,
					syscall.SO_REUSEADDR,
					1,
				)
			})
		},
	}
)

// Server sends and receives DNS messages on all available multicast
// interfaces.
type Server struct {
	wg         sync.WaitGroup
	logger     *slog.Logger
	pconn4     *ipv4.PacketConn
	pconn6     *ipv6.PacketConn
	chanPacket chan<- []byte
	chanClose  chan any
}

func (s *Server) run() {
	defer s.wg.Done()
	var (
		tickerScan = mocktime.NewTicker(30 * time.Second)
	)
	defer tickerScan.Stop()
	s.scan()
	for {
		select {
		case <-tickerScan.C:
			s.scan()
		case <-s.chanClose:
			return
		}
	}
}

// New attempts to create a new Server.
func New(cfg *Config) (*Server, error) {

	// Create the listener for IPv4 packets
	conn4, err := listenConfig.ListenPacket(
		context.Background(),
		"udp4",
		mDNSListenAddr,
	)
	if err != nil {
		return nil, err
	}

	// Create the listener for IPv6 packets
	conn6, err := listenConfig.ListenPacket(
		context.Background(),
		"udp6",
		mDNSListenAddr,
	)
	if err != nil {
		conn4.Close()
		return nil, err
	}

	// Create the server, converting the net.PacketConn to ipv4/6.PacketConn
	s := &Server{
		logger:     cfg.Logger,
		pconn4:     ipv4.NewPacketConn(conn4),
		pconn6:     ipv6.NewPacketConn(conn6),
		chanPacket: cfg.ChanPacket,
		chanClose:  make(chan any),
	}
	if s.logger == nil {
		s.logger = slog.Default()
	}
	s.logger = s.logger.With("package", "server")

	// Start the goroutines
	s.wg.Add(3)
	go s.run()
	go s.read(conn4)
	go s.read(conn6)

	return s, nil
}

// TODO: call LeaveGroup for interfaces upon shutdown

// Close shuts down the server.
func (s *Server) Close() {
	s.pconn6.Close()
	s.pconn4.Close()
	close(s.chanClose)
	s.wg.Wait()
	close(s.chanPacket)
}
