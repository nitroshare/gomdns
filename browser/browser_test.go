package browser

import (
	"testing"

	"github.com/nitroshare/gomdns/vnet"
)

func TestBrowser(t *testing.T) {
	vnet.Mock()
	defer vnet.Unmock()
	i := vnet.NewMockInterface()
	vnet.AddInterface(i)
	syncRun.Activate()
	defer syncRun.Deactivate()
	b := New(&Config{})
	defer b.Close()
	syncRun.Wait()
}
