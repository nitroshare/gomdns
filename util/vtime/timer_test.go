package vtime

import (
	"testing"
	"time"

	"github.com/nitroshare/gomdns/util/compare"
)

func TestUnmockedTimer(t *testing.T) {
	NewTimer(1 * time.Second).Stop()
}

func TestTimerStopBefore(t *testing.T) {
	Mock()
	defer Unmock()
	timer := NewTimer(1 * time.Second)
	compare.Compare(t, timer.Stop(), true, true)
}

func TestTimerStopAfter(t *testing.T) {
	Mock()
	defer Unmock()
	timer := NewTimer(1 * time.Second)
	Advance(2 * time.Second)
	<-timer.C
	compare.Compare(t, timer.Stop(), false, true)
}
