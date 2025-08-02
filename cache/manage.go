package cache

import (
	"log/slog"
	"time"

	"github.com/nitroshare/golist"
	"github.com/nitroshare/mocktime"
)

func (c *Cache) send(ch chan<- *Record, r *Record) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		select {
		case ch <- r:
		case <-c.chanClose:
		}
	}()
}

func (c *Cache) nextTrigger() <-chan time.Time {

	var (
		n           = mocktime.Now()
		nextTrigger time.Time
	)

	// Enumerate all of the records
	for e := c.entries.Front; e != nil; e = e.Next {

		// Check for elapsed triggers
		var (
			shouldQuery = false
			triggers    = e.Value.triggers
		)
		for e := triggers.Front; e != nil; e = e.Next {
			if !e.Value.After(n) {
				shouldQuery = true
				triggers.Remove(e)
			}
		}

		// If there are no triggers, the record has expired
		if triggers.Len == 0 {
			c.entries.Remove(e)
			if c.chanExpired != nil {
				r := e.Value.record
				c.logger.Debug(
					"record expired",
					slog.String("record", r.String()),
				)
				c.send(c.chanExpired, r)
			}
			continue
		}

		// Find the earliest trigger
		if nextTrigger.IsZero() || triggers.Front.Value.Before(nextTrigger) {
			nextTrigger = triggers.Front.Value
		}

		// If one of the triggers elapsed, a query is needed
		if shouldQuery && c.chanQuery != nil {
			c.send(c.chanQuery, e.Value.record)
		}
	}

	// If no records with triggers exist, return nil
	if nextTrigger.IsZero() {
		return nil
	}

	// Otherwise, return a channel that sends for the next one
	return mocktime.After(nextTrigger.Sub(n))
}

func (c *Cache) add(r *Record) {

	// Remove old records that are:
	// - of the same name/type and flush cache is set
	// - identical and TTL is set to 0
	// (Note that identical records are removed below even if TTL is set to 0;
	// this is to prevent a duplicate when the updated record is added again.)
	for e := c.entries.Front; e != nil; e = e.Next {
		var (
			sameNameType = e.Value.record.sameNameType(r)
			sameRecord   = e.Value.record.sameRecord(r)
		)
		if sameNameType && r.FlushCache || sameRecord {
			c.entries.Remove(e)
			if (r.TTL == 0 || !sameRecord) && c.chanExpired != nil {
				r := e.Value.record
				c.logger.Debug(
					"removed record",
					slog.String("record", r.String()),
				)
				c.send(c.chanExpired, r)
			}
		}
	}

	// If the record is being removed, nothing more needs to be done
	if r.TTL == 0 {
		return
	}

	// Log the new record
	c.logger.Debug(
		"added record",
		slog.String("record", r.String()),
	)

	var (
		n        = mocktime.Now()
		triggers = &golist.List[time.Time]{}
	)

	// Determine the triggers for re-querying the record (if requested)
	if c.chanQuery != nil {
		triggers.Add(n.Add(time.Duration(r.TTL) * 500 * time.Millisecond))
		triggers.Add(n.Add(time.Duration(r.TTL) * 850 * time.Millisecond))
		triggers.Add(n.Add(time.Duration(r.TTL) * 900 * time.Millisecond))
		triggers.Add(n.Add(time.Duration(r.TTL) * 950 * time.Millisecond))
	}
	triggers.Add(n.Add(time.Duration(r.TTL) * time.Second))

	// Add the entry to the list of entries
	c.entries.Add(&recordEntry{
		record:   r,
		triggers: triggers,
	})
}

func (c *Cache) lookup(name string, _type uint16) []*Record {
	records := []*Record{}
	for e := c.entries.Front; e != nil; e = e.Next {
		if e.Value.record.Name == name &&
			e.Value.record.Type == _type {
			records = append(records, e.Value.record)
		}
	}
	return records
}
