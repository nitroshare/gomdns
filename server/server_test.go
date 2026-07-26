package server

import (
	"testing"

	"github.com/nitroshare/gomdns/multicast"
	"github.com/nitroshare/gomdns/vtime"
)

func TestServerSendAll(t *testing.T) {
	vtime.Mock()
	defer vtime.Unmock()
	multicast.Mock()
	defer multicast.Unmock()
	s := New(&Config{})
	defer s.Close()
}
