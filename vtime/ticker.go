package vtime

import (
	"time"
)

// The terminology may be a bit confusing, but in this package, a Ticker is
// actually a Timer with the ticker parameter set to true. This reduces a lot
// of duplicated code in mockLoop.

// Ticker provides a drop-in replacement for time.Ticker.
type Ticker struct {
	C      <-chan time.Time
	timer  *Timer
	ticker *time.Ticker
}

// Stop ends the ticker.
func (t *Ticker) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
		return
	}
	t.timer.Stop()
}

// Reset stops the ticker and restarts it with the specified duration.
func (t *Ticker) Reset(d time.Duration) {
	if t.ticker != nil {
		t.ticker.Reset(d)
		return
	}
	loop.chanResetTimer <- &resetTimerParams{
		timer:    t.timer,
		duration: d,
	}
	<-loop.chanAny
}

// NewTicker creates a new Ticker that will send on the channel C every time
// the specified duration elapses.
func NewTicker(d time.Duration) *Ticker {
	if loop == nil {
		ticker := time.NewTicker(d)
		return &Ticker{
			C:      ticker.C,
			ticker: ticker,
		}
	}
	var (
		timer = &Timer{
			ticker:   true,
			duration: d,
		}
		t = &Ticker{
			timer: timer,
		}
	)
	loop.chanNewTimer <- t.timer
	t.C = <-loop.chanTimeChan
	return t
}
