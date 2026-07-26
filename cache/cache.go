package cache

import (
	"log/slog"
	"time"

	"github.com/nitroshare/gomdns/broadcaster"
	"github.com/nitroshare/gomdns/dns"
	"github.com/nitroshare/gomdns/list"
)

type cacheLookup struct {
	Name string
	Type uint16
}

type cacheEntry struct {
	Record   *dns.Record
	Triggers *list.List[time.Time]
}

// Cache stores records received from DNS queries and indicates when they
// should be queried again or when they expire.
type Cache struct {
	Query         *broadcaster.Broadcaster[*dns.Record]
	Expired       *broadcaster.Broadcaster[*dns.Record]
	logger        *slog.Logger
	entries       list.List[*cacheEntry]
	chanAdd       chan *dns.Record
	chanAddRet    chan any
	chanLookup    chan *cacheLookup
	chanLookupRet chan []*dns.Record
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
		case l := <-c.chanLookup:
			c.chanLookupRet <- c.lookup(l)
		case <-c.chanClose:
			return
		}
	}
}

// New returns a new Cache.
func New(cfg *Config) *Cache {
	c := &Cache{
		Query:         broadcaster.New[*dns.Record](),
		Expired:       broadcaster.New[*dns.Record](),
		logger:        cfg.Logger,
		chanAdd:       make(chan *dns.Record),
		chanAddRet:    make(chan any),
		chanLookup:    make(chan *cacheLookup),
		chanLookupRet: make(chan []*dns.Record),
		chanClose:     make(chan any),
		chanClosed:    make(chan any),
	}
	if c.logger == nil {
		c.logger = slog.Default()
	}
	c.logger = c.logger.With("package", "cache")
	go c.run()
	return c
}

// Add inserts or updates a record in the cache.
func (c *Cache) Add(record *dns.Record) {
	c.chanAdd <- record
	<-c.chanAddRet
}

// Lookup returns all records of the specified type for the provided name.
func (c *Cache) Lookup(name string, _type uint16) []*dns.Record {
	c.chanLookup <- &cacheLookup{
		Name: name,
		Type: _type,
	}
	return <-c.chanLookupRet
}

// Close shuts down the cache.
func (c *Cache) Close() {
	close(c.chanClose)
	<-c.chanClosed
	c.Query.Close()
	c.Expired.Close()
}
