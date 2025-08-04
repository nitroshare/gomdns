package server

import (
	"log/slog"
	"net"
)

func (s *Server) scan() {

	// Create a list of all interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		s.logger.Error(err.Error())
		return
	}

	for _, i := range interfaces {

		// Skip down, loopback, and non-multicast interfaces
		if i.Flags&net.FlagUp == 0 ||
			i.Flags&net.FlagLoopback != 0 ||
			i.Flags&net.FlagMulticast == 0 {
			continue
		}

		// Get the current list of joined multicast addresses
		addrs, err := i.MulticastAddrs()
		if err != nil {
			s.logger.Warn(
				err.Error(),
				slog.String("interface", i.Name),
			)
			continue
		}

		// Check to see if the mDNS addresses are already joined
		var (
			sawIPv4Addr bool
			sawIPv6Addr bool
		)
		for _, a := range addrs {
			v, ok := a.(*net.IPAddr)
			if ok && v.IP.Equal(mDNSMulticastIPv4Addr.IP) {
				sawIPv4Addr = true
				continue
			}
			if ok && v.IP.Equal(mDNSMulticastIPv6Addr.IP) {
				sawIPv6Addr = true
			}
		}

		// If not already, join the IPv4 multicast interface
		if !sawIPv4Addr {
			if err := s.pconn4.JoinGroup(&i, mDNSMulticastIPv4Addr); err != nil {
				s.logger.Warn(
					err.Error(),
					slog.String("interface", i.Name),
				)
			}
		}

		// If not already, join the IPv6 multicast interface
		if !sawIPv6Addr {
			if err := s.pconn6.JoinGroup(&i, mDNSMulticastIPv6Addr); err != nil {
				s.logger.Warn(
					err.Error(),
					slog.String("interface", i.Name),
				)
			}
		}
	}
}
