package server

import (
	"errors"
	"log/slog"
	"net"
	"sync"

	"github.com/nitroshare/gomdns/dns"
	"github.com/nitroshare/gomdns/vnet"
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

type serverConnWithAddr struct {
	Addr *net.UDPAddr
	Conn vnet.UDPConn
}

type serverListener struct {
	wg       sync.WaitGroup
	logger   *slog.Logger
	chanSend chan<- *dns.Message
	conns    []*serverConnWithAddr
}

func (l *serverListener) run(conn vnet.UDPConn) {
	defer l.wg.Done()
	for {
		b := make([]byte, 1500)
		n, addr, err := conn.ReadFrom(b)
		if err != nil {
			break
		}
		m, err := dns.ParseMessage(b[:n])
		if err == nil {
			m.Address = addr.(*net.UDPAddr).AddrPort().Addr()
			l.chanSend <- m
		}
	}
}

func newServerListener(
	logger *slog.Logger,
	i vnet.Interface,
	chanSend chan<- *dns.Message,
) (*serverListener, error) {
	flags := i.Flags()
	if flags&net.FlagUp == 0 ||
		flags&net.FlagRunning == 0 ||
		flags&net.FlagMulticast == 0 {
		return nil, errors.New("interface is not supported")
	}
	addrs, err := i.Addrs()
	if err != nil {
		return nil, err
	}
	l := &serverListener{
		logger:   logger,
		chanSend: chanSend,
	}
	for _, a := range addrs {
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
			v, err := i.Listen("udp4", IPv4Address)
			if err != nil {
				l.logger.Warn(err.Error())
				continue
			}
			l.logger.Info("listening on IPv4", "addr", ip)
			l.wg.Add(1)
			go l.run(v)
			l.conns = append(l.conns, &serverConnWithAddr{
				Addr: IPv4Address,
				Conn: v,
			})
		default:
			if ip.IsLinkLocalUnicast() {
				continue
			}
			v, err := i.Listen("udp6", IPv6Address)
			if err != nil {
				l.logger.Warn(err.Error())
				continue
			}
			l.logger.Info("listening on IPv6", "addr", ip)
			l.wg.Add(1)
			go l.run(v)
			l.conns = append(l.conns, &serverConnWithAddr{
				Addr: IPv6Address,
				Conn: v,
			})
		}
	}
	if len(l.conns) == 0 {
		return nil, errors.New("no usable addresses")
	}
	return l, nil
}

func (l *serverListener) Close() {
	for _, v := range l.conns {
		v.Conn.Close()
	}
	l.wg.Wait()
}
