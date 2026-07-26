package watcher

import (
	"testing"
	"time"

	"github.com/nitroshare/gomdns/multicast"
	"github.com/nitroshare/gomdns/vtime"
)

type watcherSet struct {
	chanAdded   <-chan multicast.Interface
	chanRemoved <-chan multicast.Interface
	watcher     *Watcher
}

func newWatcherSet() *watcherSet {
	var (
		chanAdded   = make(chan multicast.Interface)
		chanRemoved = make(chan multicast.Interface)
	)
	return &watcherSet{
		chanAdded:   chanAdded,
		chanRemoved: chanRemoved,
		watcher: NewWatcher(&Config{
			Interval:    time.Second,
			ChanAdded:   chanAdded,
			ChanRemoved: chanRemoved,
		}),
	}
}

func TestWatcher(t *testing.T) {
	vtime.Mock()
	defer vtime.Unmock()
	multicast.Mock()
	defer multicast.Unmock()
	multicast.AddMockInterface(multicast.NewMockInterface())
	s := newWatcherSet()
	defer s.watcher.Close()
	<-s.chanAdded
	multicast.ClearMockInterfaces()
	vtime.Advance(2 * time.Second)
	<-s.chanRemoved
	vtime.Advance(2 * time.Second)
	select {
	case <-s.chanAdded:
		t.Fatal("unexpected interface added")
	case <-s.chanRemoved:
		t.Fatal("unexpected interface removed")
	default:
	}
}

func TestWatcherError(t *testing.T) {
	multicast.MockWithError()
	defer multicast.Unmock()
	s := newWatcherSet()
	defer s.watcher.Close()
}

func TestCloseDuringSend(t *testing.T) {
	vtime.Mock()
	defer vtime.Unmock()
	multicast.Mock()
	defer multicast.Unmock()
	multicast.AddMockInterface(multicast.NewMockInterface())

	t.Run("chanAdded", func(t *testing.T) {
		s := newWatcherSet()
		defer s.watcher.Close()
	})

	t.Run("chanRemoved", func(t *testing.T) {
		s := newWatcherSet()
		defer s.watcher.Close()
		<-s.chanAdded
		multicast.ClearMockInterfaces()
		s.watcher.chanTest = make(chan any)
		vtime.Advance(2 * time.Second)
		<-s.watcher.chanTest
	})
}
