package server

import (
	"errors"
	"log/slog"
	"net"
	"sync"

	"github.com/nitroshare/gomdns/dns"
	"github.com/nitroshare/gomdns/multicast"
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

type serverListenerWithAddr struct {
	Addr     *net.UDPAddr
	Listener *multicast.Listener
}

type serverListener struct {
	wg        sync.WaitGroup
	logger    *slog.Logger
	chanSend  chan<- *dns.Message
	listeners []*serverListenerWithAddr
}

func (l *serverListener) run(listener *multicast.Listener) {
	defer l.wg.Done()
	for {
		p, err := listener.Read()
		if err != nil {
			break
		}
		m, err := dns.ParseMessage(p.Data)
		if err == nil {
			m.Address = p.Addr.(*net.UDPAddr).AddrPort().Addr()
			l.chanSend <- m
		}
	}
}

func newServerListener(
	logger *slog.Logger,
	i multicast.Interface,
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
			v, err := multicast.NewListener("udp4", i, IPv4Address)
			if err != nil {
				l.logger.Warn(err.Error())
				continue
			}
			l.logger.Info("listening on IPv4", "addr", ip)
			l.wg.Add(1)
			go l.run(v)
			l.listeners = append(l.listeners, &serverListenerWithAddr{
				Addr:     IPv4Address,
				Listener: v,
			})
		default:
			if ip.IsLinkLocalUnicast() {
				continue
			}
			v, err := multicast.NewListener("udp6", i, IPv6Address)
			if err != nil {
				l.logger.Warn(err.Error())
				continue
			}
			l.logger.Info("listening on IPv6", "addr", ip)
			l.wg.Add(1)
			go l.run(v)
			l.listeners = append(l.listeners, &serverListenerWithAddr{
				Addr:     IPv6Address,
				Listener: v,
			})
		}
	}
	if len(l.listeners) == 0 {
		return nil, errors.New("no usable addresses")
	}
	return l, nil
}

func (l *serverListener) Close() {
	for _, v := range l.listeners {
		v.Listener.Close()
	}
	l.wg.Wait()
}
