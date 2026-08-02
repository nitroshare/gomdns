package server

import (
	"net"

	"github.com/nitroshare/gomdns/dns"
	"github.com/nitroshare/gomdns/vnet"
)

type serverConn struct {
	conn       vnet.UDPConn
	chanSend   chan<- *dns.Message
	chanClosed chan any
}

func (c *serverConn) run() {
	defer close(c.chanClosed)
	for {
		b := make([]byte, 1500)
		n, addr, err := c.conn.ReadFrom(b)
		if err != nil {
			break
		}
		m, err := dns.ParseMessage(b[:n])
		if err == nil {
			m.Address = addr.(*net.UDPAddr).AddrPort().Addr()
			c.chanSend <- m
		}
	}
}

func newServerConn(
	i vnet.Interface,
	network string,
	addr *net.UDPAddr,
	chanSend chan<- *dns.Message,
) (*serverConn, error) {
	v, err := i.Listen(network, addr)
	if err != nil {
		return nil, err
	}
	c := &serverConn{
		conn:       v,
		chanSend:   chanSend,
		chanClosed: make(chan any),
	}
	go c.run()
	return c, nil
}

func (c *serverConn) Close() {
	c.conn.Close()
	<-c.chanClosed
}
