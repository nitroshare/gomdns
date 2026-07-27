package watcher

import (
	"errors"
	"testing"
	"time"

	"github.com/nitroshare/gomdns/vnet"
	"github.com/nitroshare/gomdns/vtime"
)

type watcherSet struct {
	chanAdded   <-chan vnet.Interface
	chanRemoved <-chan vnet.Interface
	watcher     *Watcher
}

func newWatcherSet() *watcherSet {
	var (
		chanAdded   = make(chan vnet.Interface)
		chanRemoved = make(chan vnet.Interface)
	)
	return &watcherSet{
		chanAdded:   chanAdded,
		chanRemoved: chanRemoved,
		watcher: New(&Config{
			Interval:    time.Second,
			ChanAdded:   chanAdded,
			ChanRemoved: chanRemoved,
		}),
	}
}

func TestWatcher(t *testing.T) {
	vtime.Mock()
	defer vtime.Unmock()
	vnet.Mock()
	defer vnet.Unmock()
	vnet.AddInterface(vnet.NewMockInterface())
	s := newWatcherSet()
	defer s.watcher.Close()
	<-s.chanAdded
	vnet.ClearInterfaces()
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
	vnet.Interfaces = func() ([]vnet.Interface, error) {
		return nil, errors.New("test")
	}
	defer vnet.Unmock()
	s := newWatcherSet()
	defer s.watcher.Close()
}

func TestCloseDuringSend(t *testing.T) {
	vtime.Mock()
	defer vtime.Unmock()
	vnet.Mock()
	defer vnet.Unmock()
	vnet.AddInterface(vnet.NewMockInterface())

	t.Run("chanAdded", func(t *testing.T) {
		s := newWatcherSet()
		defer s.watcher.Close()
	})

	t.Run("chanRemoved", func(t *testing.T) {
		s := newWatcherSet()
		defer s.watcher.Close()
		<-s.chanAdded
		vnet.ClearInterfaces()
		s.watcher.chanTest = make(chan any)
		vtime.Advance(2 * time.Second)
		<-s.watcher.chanTest
	})
}
