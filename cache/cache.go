package cache

import (
	"log/slog"
	"sync"
	"time"

	"github.com/nitroshare/gomdns/util/list"
)

var (
	fnNow   = time.Now
	fnAfter = time.After
)

type recordEntry struct {
	record   *Record
	triggers *list.List[time.Time]
}

// Cache stores records received from DNS queries and sends on the
// shouldQuery channel when records are about to expire.
type Cache struct {
	mutex       sync.Mutex
	logger      *slog.Logger
	entries     *list.List[*recordEntry]
	chanQuery   chan<- *Record
	chanExpired chan<- *Record
	chanAdd     chan *Record
	chanClosed  chan any
}

func (c *Cache) run() {
	defer close(c.chanClosed)
	for {
		select {
		case <-c.nextTrigger():
		case r, ok := <-c.chanAdd:
			if !ok {
				return
			}
			c.add(r)
		}
	}
}

// New returns a new Cache instance.
func New(cfg *Config) *Cache {
	c := &Cache{
		logger:      cfg.Logger,
		entries:     &list.List[*recordEntry]{},
		chanQuery:   cfg.ChanQuery,
		chanExpired: cfg.ChanExpired,
		chanAdd:     make(chan *Record),
		chanClosed:  make(chan any),
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
}

// LookupByName returns all records for the provided name.
func (c *Cache) LookupByName(name string) []*Record {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	records := []*Record{}
	for e := c.entries.Front; e != nil; e = e.Next {
		if e.Value.record.Name == name {
			records = append(records, e.Value.record)
		}
	}
	return records
}

// LookupByNameAndType returns all records of the specified type for the
// provided name.
func (c *Cache) LookupByNameAndType(name string, _type uint16) []*Record {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	records := []*Record{}
	for e := c.entries.Front; e != nil; e = e.Next {
		if e.Value.record.Name == name &&
			e.Value.record.Type == _type {
			records = append(records, e.Value.record)
		}
	}
	return records
}

// Close shuts down the cache.
func (c *Cache) Close() {
	close(c.chanAdd)
	<-c.chanClosed
}
