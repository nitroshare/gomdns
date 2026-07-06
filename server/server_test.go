package server

import (
	"testing"

	"github.com/nitroshare/gomulticast"
	"github.com/nitroshare/gotime"
)

func TestServer(t *testing.T) {
	gotime.Mock()
	defer gotime.Unmock()
	gomulticast.Mock()
	defer gomulticast.Unmock()
	s := New()
	defer s.Close()
}
