package multicast

import (
	"errors"
	"net"
	"testing"

	"github.com/nitroshare/compare"
)

var (
	testUDPAddr = &net.UDPAddr{
		IP:   net.IP([]byte{1, 2, 3, 4}),
		Port: 1234,
	}
	testUDPData = []byte("test")
	testPacket  = &Packet{
		Addr: testUDPAddr,
		Data: testUDPData,
	}
	testError = errors.New("test")
)

func TestNewListenerError(t *testing.T) {
	origNetListenMulticastUDP = func(_ string, _ *net.Interface, _ *net.UDPAddr) (*net.UDPConn, error) {
		return nil, testError
	}
	defer func() {
		origNetListenMulticastUDP = net.ListenMulticastUDP
	}()
	_, err := NewListener("", NewMockInterface(), nil)
	compare.Compare(t, err, testError, true)
}

func TestListenerRead(t *testing.T) {
	Mock()
	defer Unmock()
	var (
		i    = NewMockInterface()
		l, _ = NewListener("udp4", i, testUDPAddr)
	)
	defer l.Close()
	t.Run("test successful read", func(t *testing.T) {
		i.QueueForRead(testPacket)
		p, err := l.Read()
		compare.Compare(t, p.Addr, net.Addr(testUDPAddr), true)
		compare.Compare(t, string(p.Data), string(testUDPData), true)
		compare.Compare(t, err, nil, true)
	})
	t.Run("test unsuccessful read", func(t *testing.T) {
		i.errOnNextRead = true
		p, err := l.Read()
		compare.Compare(t, p, nil, true)
		compare.Compare(t, err, net.ErrClosed, true)
	})
}

func TestListenerWrite(t *testing.T) {
	Mock()
	defer Unmock()
	var (
		i    = NewMockInterface()
		l, _ = NewListener("udp4", i, testUDPAddr)
	)
	defer l.Close()
	l.Write(testPacket)
	p, err := i.DequeueWrite()
	compare.Compare(t, p.Addr, net.Addr(testUDPAddr), true)
	compare.Compare(t, string(p.Data), string(testUDPData), true)
	compare.Compare(t, err, nil, true)
}
