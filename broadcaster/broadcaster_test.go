package broadcaster

import (
	"testing"

	"github.com/nitroshare/gomdns/compare"
)

func TestBroadcaster(t *testing.T) {
	var (
		b = New[int]()
		s = b.Subscribe()
	)
	b.Send(5)
	compare.Compare(t, <-s, 5, true)
	b.Unsubscribe(s)
	_, ok := <-s
	compare.Compare(t, ok, false, true)
	s = b.Subscribe()
	b.Close()
}
