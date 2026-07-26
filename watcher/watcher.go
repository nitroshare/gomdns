package watcher

import (
	"time"

	"github.com/nitroshare/gomdns/multicast"
	"github.com/nitroshare/gotime"
)

// Watcher monitors available network interfaces and notifies when one is
// added or removed.
type Watcher struct {
	chanTest    chan any
	chanAdded   chan<- multicast.Interface
	chanRemoved chan<- multicast.Interface
	chanClose   chan any
	chanClosed  chan any
}

func (w *Watcher) diff(m map[string]multicast.Interface) map[string]multicast.Interface {
	interfaces, err := multicast.Interfaces()
	if err != nil {
		return m
	}
	m2 := map[string]multicast.Interface{}
	for _, i := range interfaces {
		m2[i.Interface().Name] = i
	}
	for k, v := range m {
		if _, ok := m2[k]; !ok {
			select {
			case w.chanRemoved <- v:
			case <-w.chanClose:
				return nil
			}
		}
	}
	for k, v := range m2 {
		if _, ok := m[k]; !ok {
			select {
			case w.chanAdded <- v:
			case <-w.chanClose:
				return nil
			}
		}
	}
	return m2
}

func (w *Watcher) run(interval time.Duration) {
	defer close(w.chanClosed)
	defer close(w.chanRemoved)
	defer close(w.chanAdded)
	t := gotime.NewTicker(interval)
	defer t.Stop()
	m := w.diff(map[string]multicast.Interface{})
	for {
		select {
		case <-t.C:
			if w.chanTest != nil {
				close(w.chanTest)
			}
			m = w.diff(m)
			if m == nil {
				return
			}
		case <-w.chanClose:
			return
		}
	}
}

// NewWatcher creates a new Watcher instance.
func NewWatcher(cfg *Config) *Watcher {
	w := &Watcher{
		chanAdded:   cfg.ChanAdded,
		chanRemoved: cfg.ChanRemoved,
		chanClose:   make(chan any),
		chanClosed:  make(chan any),
	}
	go w.run(cfg.Interval)
	return w
}

// Close shuts down the watcher.
func (w *Watcher) Close() {
	close(w.chanClose)
	<-w.chanClosed
}
