package gomdns

import (
	"net"
)

var (
	// IPv4Address is the IPv4 address for mDNS datagrams.
	IPv4Address = &net.UDPAddr{
		IP:   net.IPv4(224, 0, 0, 251),
		Port: 5353,
	}

	// IPv6Address is the IPv6 address for mDNS datagrams.
	IPv6Address = &net.UDPAddr{
		IP:   net.ParseIP("ff02::fb"),
		Port: 5353,
	}
)
