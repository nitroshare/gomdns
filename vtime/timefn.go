package vtime

import (
	"time"
)

// Set explicitly sets the mocked time.
func Set(t time.Time) {
	loop.chanSet <- t
	<-loop.chanAny
}

// Advance advances the mocked time by the specified duration.
func Advance(d time.Duration) {
	loop.chanAdvance <- d
	<-loop.chanAny
}

// AdvanceToAfter is a synchronization function that will advance to the time
// of the next After() call, waiting for a call to After() if necessary. Note
// that this will also include timers.
func AdvanceToAfter() {
	loop.chanAdvanceToAfter <- nil
	<-loop.chanAny
}
