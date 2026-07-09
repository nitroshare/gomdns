package browser

import (
	"testing"

	"github.com/nitroshare/gomulticast"
)

func TestBrowser(t *testing.T) {
	gomulticast.Mock()
	defer gomulticast.Unmock()
	b := New(&Config{})
	defer b.Close()
}
