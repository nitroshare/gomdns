package cache

import (
	"log/slog"
	"time"

	"github.com/nitroshare/gomdns/dns"
	"github.com/nitroshare/gomdns/list"
	"github.com/nitroshare/gomdns/vtime"
)

func (c *Cache) nextTrigger() <-chan time.Time {
	var (
		n           = vtime.Now()
		nextTrigger time.Time
	)
	for e := c.entries.Front; e != nil; e = e.Next {
		var (
			shouldQuery = false
			triggers    = e.Value.Triggers
		)
		for e := triggers.Front; e != nil; e = e.Next {
			if !e.Value.After(n) {
				shouldQuery = true
				triggers.Remove(e)
			}
		}
		if triggers.Len == 0 {
			c.entries.Remove(e)
			r := e.Value.Record
			c.logger.Debug(
				"record expired",
				slog.String("record", r.String()),
			)
			c.Expired.Send(r)
			continue
		}
		if nextTrigger.IsZero() || triggers.Front.Value.Before(nextTrigger) {
			nextTrigger = triggers.Front.Value
		}
		if shouldQuery {
			c.Query.Send(e.Value.Record)
		}
	}
	if nextTrigger.IsZero() {
		return nil
	}
	return vtime.After(nextTrigger.Sub(n))
}

func (c *Cache) add(r *dns.Record) {
	for e := c.entries.Front; e != nil; e = e.Next {
		var (
			sameNameType = e.Value.Record.SameNameAndType(r)
			sameRecord   = e.Value.Record.SameRecord(r)
		)
		if sameNameType && r.FlushCache || sameRecord {
			c.entries.Remove(e)
			if r.Ttl == 0 || !sameRecord {
				r := e.Value.Record
				c.logger.Debug(
					"removed record",
					slog.String("record", r.String()),
				)
				c.Expired.Send(r)
			}
		}
	}
	if r.Ttl == 0 {
		return
	}
	c.logger.Debug(
		"added / updated record",
		slog.String("record", r.String()),
	)
	var (
		n        = vtime.Now()
		triggers = &list.List[time.Time]{}
	)
	triggers.Add(n.Add(time.Duration(r.Ttl) * 500 * time.Millisecond))
	triggers.Add(n.Add(time.Duration(r.Ttl) * 850 * time.Millisecond))
	triggers.Add(n.Add(time.Duration(r.Ttl) * 900 * time.Millisecond))
	triggers.Add(n.Add(time.Duration(r.Ttl) * 950 * time.Millisecond))
	triggers.Add(n.Add(time.Duration(r.Ttl) * time.Second))
	c.entries.Add(&cacheEntry{
		Record:   r,
		Triggers: triggers,
	})
}

func (c *Cache) lookup(l *cacheLookup) []*dns.Record {
	records := []*dns.Record{}
	for e := c.entries.Front; e != nil; e = e.Next {
		if e.Value.Record.Name == l.Name &&
			e.Value.Record.Type == l.Type {
			records = append(records, e.Value.Record)
		}
	}
	return records
}
