package server

import (
	"errors"
	"net"
	"net/netip"
	"testing"

	"github.com/nitroshare/gomdns/compare"
	"github.com/nitroshare/gomdns/dns"
	"github.com/nitroshare/gomdns/vnet"
)

var (
	testAddr    = netip.MustParseAddr("1.2.3.4")
	testUDPAddr = &net.UDPAddr{
		IP: testAddr.AsSlice(),
	}
)

func TestServerConn(t *testing.T) {
	var (
		i        = vnet.NewMockInterface()
		chanSend = make(chan *dns.Message)
	)
	c, err := newServerConn(i, "udp4", nil, chanSend)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	i.EnqueueRead(&vnet.Packet{
		Addr: testUDPAddr,
		Data: []byte{
			0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		},
	})
	m := <-chanSend
	compare.Compare(t, m.Address, testAddr, true)
	compare.Compare(t, m.TransactionID, 1, true)
}

type testServerInterface struct{}

func (i *testServerInterface) Name() string               { return "" }
func (i *testServerInterface) Addrs() ([]net.Addr, error) { return nil, nil }
func (i *testServerInterface) Flags() net.Flags           { return 0 }
func (i *testServerInterface) Listen(string, *net.UDPAddr) (vnet.UDPConn, error) {
	return nil, errors.New("test")
}

func TestServerConnError(t *testing.T) {
	_, err := newServerConn(&testServerInterface{}, "", nil, nil)
	compare.Compare(t, err != nil, true, true)
}
