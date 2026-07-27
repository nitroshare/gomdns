package browser

import (
	"testing"

	"github.com/nitroshare/gomdns/multicast"
)

func TestBrowser(t *testing.T) {
	multicast.Mock()
	defer multicast.Unmock()
	i := multicast.NewMockInterface()
	multicast.AddMockInterface(i)
	syncRun.Activate()
	defer syncRun.Deactivate()
	b := New(&Config{})
	defer b.Close()
	syncRun.Wait()
}
