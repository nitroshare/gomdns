package cache

import (
	"time"

	"github.com/nitroshare/golist"
	"github.com/nitroshare/mocktime"
)

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
				c.chanExpired <- e.Value.record
			}
			continue
		}

		// Find the earliest trigger
		if nextTrigger.IsZero() || triggers.Front.Value.Before(nextTrigger) {
			nextTrigger = triggers.Front.Value
		}

		// If one of the triggers elapsed, a query is needed
		if shouldQuery && c.chanQuery != nil {
			c.chanQuery <- e.Value.record
		}
	}

	// If no records with triggers exist, return nil
	if nextTrigger.IsZero() {
		return nil
	}

	// Otherwise, return a channel that sends for the next one
	return mocktime.After(nextTrigger.Sub(n))
}

func (c *Cache) add(record *Record) {

	// If records with the same name / type exist, keep them only if they are
	// different or the flush cache bit is not set
	for e := c.entries.Front; e != nil; e = e.Next {
		if e.Value.record.sameNameType(record) && record.FlushCache ||
			e.Value.record.sameRecord(record) {
			c.entries.Remove(e)
			if record.TTL == 0 && c.chanExpired != nil {
				c.chanExpired <- e.Value.record
			}
		}
	}

	// If the record is being removed, nothing more needs to be done
	if record.TTL == 0 {
		return
	}

	var (
		n        = mocktime.Now()
		triggers = &golist.List[time.Time]{}
	)

	// Determine the triggers for re-querying the record
	triggers.Add(n.Add(time.Duration(record.TTL) * 500 * time.Millisecond))
	triggers.Add(n.Add(time.Duration(record.TTL) * 850 * time.Millisecond))
	triggers.Add(n.Add(time.Duration(record.TTL) * 900 * time.Millisecond))
	triggers.Add(n.Add(time.Duration(record.TTL) * 950 * time.Millisecond))
	triggers.Add(n.Add(time.Duration(record.TTL) * time.Second))

	// Add the entry to the list of entries
	c.entries.Add(&recordEntry{
		record:   record,
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
