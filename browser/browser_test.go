package browser

import (
	"testing"

	"github.com/nitroshare/gomdns/multicast"
)

func TestBrowser(t *testing.T) {
	multicast.Mock()
	defer multicast.Unmock()
	b := New(&Config{})
	defer b.Close()
}
