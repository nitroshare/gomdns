package server

import (
	"errors"
	"net"
)

func (s *Server) scan() {
	interfaces, err := net.Interfaces()
	if err != nil {
		s.logger.Error(err.Error())
		return
	}
	for _, i := range interfaces {
		if i.Flags&net.FlagUp == 0 ||
			i.Flags&net.FlagLoopback != 0 ||
			i.Flags&net.FlagMulticast == 0 {
			continue
		}

		// No error checking is done here because it would result in a lot of
		// false positives. For example, joining a multicast group twice
		// throws an error on Windows. Interface.MulticastAddrs seems like a
		// solution but it returns addresses that aren't joined yet.

		s.pconn4.JoinGroup(&i, mDNSMulticastIPv4Addr)
		s.pconn6.JoinGroup(&i, mDNSMulticastIPv6Addr)
	}
}

func (s *Server) read(conn net.PacketConn) {
	defer s.wg.Done()
	b := make([]byte, 1500)
	for {
		n, _, err := conn.ReadFrom(b)
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				s.logger.Error(err.Error())
			}
			break
		}
		s.chanPacket <- b[:n]
	}
}
