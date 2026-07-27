package server

import (
	"errors"
	"net"
	"strings"
	"testing"

	"github.com/nitroshare/gomdns/compare"
	"github.com/nitroshare/gomdns/dns"
	"github.com/nitroshare/gomdns/multicast"
)

const (
	validFlags = net.FlagUp | net.FlagRunning | net.FlagMulticast
)

func TestServer(t *testing.T) {
	multicast.Mock()
	defer multicast.Unmock()
	syncAdd.Activate()
	defer syncAdd.Deactivate()
	syncAddSuccess.Activate()
	defer syncAddSuccess.Deactivate()
	for _, v := range []struct {
		Name           string
		MockAddrs      []net.Addr
		MockAddrsError error
		MockFlags      net.Flags
		Success        bool
	}{
		{
			Name:           "Addrs error",
			MockAddrsError: errors.New("error"),
			MockFlags:      validFlags,
		},
		{
			Name: "Bad address",
			MockAddrs: []net.Addr{
				nil,
			},
			MockFlags: validFlags,
		},
		{
			Name: "Missing flags",
			MockAddrs: []net.Addr{
				&net.IPNet{IP: net.ParseIP("1.2.3.4")},
			},
		},
		{
			Name: "Loopback (IPv4)",
			MockAddrs: []net.Addr{
				&net.IPNet{IP: net.ParseIP("127.0.0.1")},
			},
			MockFlags: validFlags,
		},
		{
			Name: "Loopback (IPv6)",
			MockAddrs: []net.Addr{
				&net.IPNet{IP: net.ParseIP("fe80::")},
			},
			MockFlags: validFlags,
		},
		{
			Name: "Valid IPv4",
			MockAddrs: []net.Addr{
				&net.IPNet{IP: net.ParseIP("1.2.3.4")},
			},
			MockFlags: validFlags,
			Success:   true,
		},
		{
			Name: "Valid IPv6",
			MockAddrs: []net.Addr{
				&net.IPNet{IP: net.ParseIP("::1")},
			},
			MockFlags: validFlags,
			Success:   true,
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			i := multicast.NewMockInterface()
			i.MockAddrs = v.MockAddrs
			i.MockAddrsError = v.MockAddrsError
			i.MockFlags = v.MockFlags
			multicast.AddMockInterface(i)
			defer multicast.ClearMockInterfaces()
			s := New(&Config{})
			defer s.Close()
			syncAdd.Wait()
			if v.Success {
				syncAddSuccess.Wait()
			}
		})
	}
}

func TestServerSendReceive(t *testing.T) {
	multicast.Mock()
	defer multicast.Unmock()
	syncAdd.Activate()
	defer syncAdd.Deactivate()
	i := multicast.NewMockInterface()
	i.MockAddrs = []net.Addr{
		&net.IPNet{IP: net.ParseIP("1.2.3.4")},
	}
	i.MockFlags = validFlags
	multicast.AddMockInterface(i)
	defer multicast.ClearMockInterfaces()
	s := New(&Config{})
	defer s.Close()
	syncAdd.Wait()
	s.Send(&dns.Message{})
	p, err := i.DequeueWrite()
	if err != nil {
		t.Fatal(err)
	}
	compare.Compare(t, p, nil, false)
	chanReceived := s.Received.Subscribe()
	b, err := (&dns.Message{}).Serialize()
	if err != nil {
		t.Fatal(err)
	}
	i.QueueForRead(&multicast.Packet{
		Addr: &net.UDPAddr{IP: net.ParseIP("1.2.3.4")},
		Data: b,
	})
	<-chanReceived
	s.Send(&dns.Message{
		Questions: []*dns.Question{
			{
				Name: strings.Repeat("0", 64),
			},
		},
	})
}
