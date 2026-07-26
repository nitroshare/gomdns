package vtime

import (
	"testing"
	"time"

	"github.com/nitroshare/compare"
)

func TestAfter(t *testing.T) {
	Mock()
	defer Unmock()
	var (
		n          = Now()
		d          = 1 * time.Second
		chanClosed = make(chan any)
	)

	// Set this to non-nil to force the execution order of AdvanceToAfter and
	// the call to After in the secondary goroutine
	loop.chanTest = make(chan any)

	go func() {
		defer close(chanClosed)
		<-After(d)
	}()
	AdvanceToAfter()
	compare.Compare(t, Now(), n.Add(d), true)
	<-chanClosed
}
