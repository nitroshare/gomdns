package browser

import (
	"testing"

	"github.com/nitroshare/gomdns/server"
	"github.com/nitroshare/gomulticast"
)

func TestBrowser(t *testing.T) {
	gomulticast.Mock()
	defer gomulticast.Unmock()
	var (
		s = server.New()
		b = New(&Config{
			Server: s,
		})
	)
	defer s.Close()
	defer b.Close()
}
