package multicast

import (
	"net"
)

var (
	origNetListenMulticastUDP = net.ListenMulticastUDP
)

type udpConn interface {
	ReadFrom([]byte) (int, net.Addr, error)
	WriteTo([]byte, net.Addr) (int, error)
	Close() error
}

var (
	listenMulticastUDP func(string, Interface, *net.UDPAddr) (udpConn, error)
)

func netListenMulticastUDP(
	network string,
	i Interface,
	addr *net.UDPAddr,
) (udpConn, error) {
	return origNetListenMulticastUDP(network, i.Interface(), addr)
}

func mockListenMulticastUDP(
	network string,
	i Interface,
	addr *net.UDPAddr,
) (udpConn, error) {
	return &mockUDPConn{i.(*MockInterface)}, nil
}

// Listener provides a simple interface for sending and receiving UDP packets
// on a network interface.
type Listener struct {
	conn udpConn
	mtu  int
}

// NewListener creates a new Listener for the provided network, interface, and
// multicast address. Despite what the name suggests, it can be used for
// sending packets as well.
func NewListener(network string, i Interface, addr *net.UDPAddr) (*Listener, error) {
	c, err := listenMulticastUDP(
		network,
		i,
		addr,
	)
	if err != nil {
		return nil, err
	}
	return &Listener{
		conn: c,
		mtu:  i.Interface().MTU,
	}, nil
}

// Read is a wrapper over the ReadFrom method of the underlying connection.
// Note that a packet may be returned with an error. The packet should be
// considered before handling the error.
func (l *Listener) Read() (*Packet, error) {
	b := make([]byte, l.mtu)
	n, addr, err := l.conn.ReadFrom(b)
	if addr == nil {
		return nil, err
	}
	return &Packet{
		Addr: addr,
		Data: b[:n],
	}, err
}

// Write is a wrapper over the WriteTo method of the underlying connection.
func (l *Listener) Write(p *Packet) (int, error) {
	return l.conn.WriteTo(p.Data, p.Addr)
}

// Close shuts down the listener by closing the underlying connection.
func (l *Listener) Close() {
	l.conn.Close()
}
