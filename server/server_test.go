package server

import (
	"errors"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/nitroshare/gomdns/compare"
	"github.com/nitroshare/gomdns/dns"
	"github.com/nitroshare/gomdns/vnet"
	"github.com/nitroshare/gomdns/vtime"
)

const (
	validFlags = net.FlagUp | net.FlagRunning | net.FlagMulticast
)

func TestServer(t *testing.T) {
	vnet.Mock()
	defer vnet.Unmock()
	syncAdd.Activate()
	defer syncAdd.Deactivate()
	syncAddSuccess.Activate()
	defer syncAddSuccess.Deactivate()
	for _, v := range []struct {
		Name            string
		MockAddrs       []net.Addr
		MockAddrsError  error
		MockFlags       net.Flags
		MockListenError error
		Success         bool
	}{
		{
			Name:           "Addrs error",
			MockAddrsError: errors.New("error"),
			MockFlags:      validFlags,
		},
		{
			Name: "Missing flags",
			MockAddrs: []net.Addr{
				&net.IPNet{IP: net.ParseIP("1.2.3.4")},
			},
		},
		{
			Name: "Bad address",
			MockAddrs: []net.Addr{
				nil,
			},
			MockFlags: validFlags,
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
		{
			Name: "Valid IPv4 but Listen fails",
			MockAddrs: []net.Addr{
				&net.IPNet{IP: net.ParseIP("1.2.3.4")},
			},
			MockFlags:       validFlags,
			MockListenError: errors.New("test"),
		},
		{
			Name: "Valid IPv6 but Listen fails",
			MockAddrs: []net.Addr{
				&net.IPNet{IP: net.ParseIP("::1")},
			},
			MockFlags:       validFlags,
			MockListenError: errors.New("test"),
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			i := vnet.NewMockInterface()
			i.MockAddrs = v.MockAddrs
			i.MockAddrsError = v.MockAddrsError
			i.MockFlags = v.MockFlags
			i.MockListenError = v.MockListenError
			vnet.AddInterface(i)
			defer vnet.ClearInterfaces()
			s := New(&Config{})
			defer s.Close()
			syncAdd.Wait()
			if v.Success {
				syncAddSuccess.Wait()
			}
		})
	}
}

func TestServerInterfaceRemoved(t *testing.T) {
	vtime.Mock()
	defer vtime.Unmock()
	vnet.Mock()
	defer vnet.Unmock()
	syncAddSuccess.Activate()
	defer syncAddSuccess.Deactivate()
	syncRemoveSuccess.Activate()
	defer syncRemoveSuccess.Deactivate()
	i := vnet.NewMockInterface()
	i.MockAddrs = []net.Addr{
		&net.IPNet{IP: net.ParseIP("1.2.3.4")},
	}
	i.MockFlags = validFlags
	vnet.AddInterface(i)
	s := New(&Config{})
	defer s.Close()
	syncAddSuccess.Wait()
	vnet.ClearInterfaces()
	vtime.Advance(1 * time.Minute)
	syncRemoveSuccess.Wait()
}

func TestServerSendReceive(t *testing.T) {
	vnet.Mock()
	defer vnet.Unmock()
	syncAddSuccess.Activate()
	defer syncAddSuccess.Deactivate()
	i := vnet.NewMockInterface()
	i.MockAddrs = []net.Addr{
		&net.IPNet{IP: net.ParseIP("1.2.3.4")},
	}
	i.MockFlags = validFlags
	vnet.AddInterface(i)
	s := New(&Config{})
	defer s.Close()
	syncAddSuccess.Wait()
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
	i.EnqueueRead(&vnet.Packet{
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
