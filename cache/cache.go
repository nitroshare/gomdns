package cache

import (
	"log/slog"
	"sync"
	"time"

	"github.com/nitroshare/golist"
)

type recordEntry struct {
	record   *Record
	triggers *golist.List[time.Time]
}

type lookupParams struct {
	name  string
	_type uint16
}

// Cache stores records received from DNS queries and sends on the
// shouldQuery channel when records are about to expire.
type Cache struct {
	wg            sync.WaitGroup
	logger        *slog.Logger
	entries       *golist.List[*recordEntry]
	chanQuery     chan<- *Record
	chanExpired   chan<- *Record
	chanAdd       chan *Record
	chanAddRet    chan any
	chanLookup    chan *lookupParams
	chanLookupRet chan []*Record
	chanClose     chan any
	chanClosed    chan any
}

func (c *Cache) run() {
	defer close(c.chanClosed)
	for {
		select {
		case <-c.nextTrigger():
		case r := <-c.chanAdd:
			c.add(r)
			c.chanAddRet <- nil
		case p := <-c.chanLookup:
			c.chanLookupRet <- c.lookup(p.name, p._type)
		case <-c.chanClose:
			return
		}
	}
}

// New returns a new Cache instance.
func New(cfg *Config) *Cache {
	c := &Cache{
		logger:        cfg.Logger,
		entries:       &golist.List[*recordEntry]{},
		chanQuery:     cfg.ChanQuery,
		chanExpired:   cfg.ChanExpired,
		chanAdd:       make(chan *Record),
		chanAddRet:    make(chan any),
		chanLookup:    make(chan *lookupParams),
		chanLookupRet: make(chan []*Record),
		chanClose:     make(chan any),
		chanClosed:    make(chan any),
	}
	if c.logger == nil {
		c.logger = slog.Default()
	}
	go c.run()
	return c
}

// Add adds a record to the cache.
func (c *Cache) Add(record *Record) {
	c.chanAdd <- record
	<-c.chanAddRet
}

// Lookup returns all records of the specified type for the provided name.
func (c *Cache) Lookup(name string, _type uint16) []*Record {
	c.chanLookup <- &lookupParams{
		name:  name,
		_type: _type,
	}
	return <-c.chanLookupRet
}

// Close shuts down the cache.
func (c *Cache) Close() {
	close(c.chanClose)
	c.wg.Wait()
	<-c.chanClosed
	if c.chanQuery != nil {
		close(c.chanQuery)
	}
	if c.chanExpired != nil {
		close(c.chanExpired)
	}
}
