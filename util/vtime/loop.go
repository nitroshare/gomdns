package vtime

import (
	"time"

	"github.com/nitroshare/golist"
)

type afterChanData struct {
	expiry time.Time
	ch     chan time.Time
	timer  *Timer
}

type resetTimerParams struct {
	timer    *Timer
	duration time.Duration
}

type mockLoop struct {
	mockTime           time.Time
	waitingForAfter    bool
	afterChan          golist.List[*afterChanData]
	chanNow            chan any
	chanAfter          chan time.Duration
	chanSet            chan time.Time
	chanAdvance        chan time.Duration
	chanAdvanceToAfter chan any
	chanNewTimer       chan *Timer
	chanStopTimer      chan *Timer
	chanResetTimer     chan *resetTimerParams
	chanAny            chan any
	chanBool           chan bool
	chanTime           chan time.Time
	chanTimeChan       chan (<-chan time.Time)
	chanTest           chan any
	chanClose          chan any
}

func (m *mockLoop) send(expiry time.Time, ch chan<- time.Time) {
	go func() {
		ch <- expiry
	}()
}

func (m *mockLoop) after(d time.Duration, t *Timer) <-chan time.Time {
	var (
		expiry = m.mockTime.Add(d)
		ch     = make(chan time.Time)
	)
	if m.waitingForAfter {
		m.mockTime = expiry
		m.send(expiry, ch)
		m.waitingForAfter = false
		m.chanAny <- nil
	} else {
		m.afterChan.Add(&afterChanData{
			expiry: expiry,
			ch:     ch,
			timer:  t,
		})
	}
	return ch
}

func (m *mockLoop) sendElapsed() {
	for e := m.afterChan.Front; e != nil; e = e.Next {
		d := e.Value
		if !d.expiry.After(m.mockTime) {
			m.send(d.expiry, d.ch)
			if d.timer != nil && d.timer.ticker {
				d.expiry = d.expiry.Add(d.timer.duration)
			} else {
				m.afterChan.Remove(e)
			}
		}
	}
}

func (m *mockLoop) findEarliest() time.Time {
	var earliestTime time.Time
	for e := m.afterChan.Front; e != nil; e = e.Next {
		if earliestTime.IsZero() || e.Value.expiry.Before(earliestTime) {
			earliestTime = e.Value.expiry
		}
	}
	return earliestTime
}

func (m *mockLoop) stopTimer(t *Timer) bool {
	for e := m.afterChan.Front; e != nil; e = e.Next {
		if t == e.Value.timer {
			m.afterChan.Remove(e)
			return true
		}
	}
	return false
}

func (m *mockLoop) resetTimer(t *Timer, d time.Duration) {
	for e := m.afterChan.Front; e != nil; e = e.Next {
		if t == e.Value.timer {
			e.Value.expiry = m.mockTime.Add(d)
			t.duration = d
			break
		}
	}
}

func (m *mockLoop) run() {
	defer close(m.chanClose)
	for {
		select {
		case <-m.chanNow:
			m.chanTime <- m.mockTime
		case d := <-m.chanAfter:
			m.chanTimeChan <- m.after(d, nil)
		case t := <-m.chanSet:
			m.mockTime = t
			m.sendElapsed()
			m.chanAny <- nil
		case d := <-m.chanAdvance:
			m.mockTime = m.mockTime.Add(d)
			m.sendElapsed()
			m.chanAny <- nil
		case <-m.chanAdvanceToAfter:
			earliestTime := m.findEarliest()
			if earliestTime.IsZero() {
				if m.chanTest != nil {
					m.chanTest <- nil
				}
				m.waitingForAfter = true
				continue
			}
			m.mockTime = earliestTime
			m.sendElapsed()
			m.chanAny <- nil
		case t := <-m.chanNewTimer:
			m.chanTimeChan <- m.after(t.duration, t)
		case t := <-m.chanStopTimer:
			m.chanBool <- m.stopTimer(t)
		case d := <-m.chanResetTimer:
			m.resetTimer(d.timer, d.duration)
			m.chanAny <- nil
		case <-m.chanClose:
			return
		}
	}
}

func newMockLoop() *mockLoop {
	m := &mockLoop{
		afterChan:          golist.List[*afterChanData]{},
		chanNow:            make(chan any),
		chanAfter:          make(chan time.Duration),
		chanSet:            make(chan time.Time),
		chanAdvance:        make(chan time.Duration),
		chanAdvanceToAfter: make(chan any),
		chanNewTimer:       make(chan *Timer),
		chanStopTimer:      make(chan *Timer),
		chanResetTimer:     make(chan *resetTimerParams),
		chanAny:            make(chan any),
		chanBool:           make(chan bool),
		chanTime:           make(chan time.Time),
		chanTimeChan:       make(chan (<-chan time.Time)),
		chanClose:          make(chan any),
	}
	go m.run()
	return m
}

func (m *mockLoop) close() {
	m.chanClose <- nil
	<-m.chanClose
}

var loop *mockLoop
