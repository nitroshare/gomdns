package cache

import (
	"maps"
)

// Record represents a structured DNS record.
type Record struct {

	// Name is the name of the record.
	Name string

	// Type is the type of the record (eg., A, AAAA, SRV).
	Type uint16

	// Target is the domain the record points to.
	Target string

	// Attributes contains the key=value from TXT records.
	Attributes map[string]string

	// FlushCache indicates that this record should supercede all others of
	// the same name and type.
	FlushCache bool

	// TTL is the time-to-live for the record.
	TTL int
}

func (r *Record) sameNameType(record *Record) bool {
	return r.Name == record.Name && r.Type == record.Type
}

func (r *Record) sameRecord(record *Record) bool {
	return r.sameNameType(record) &&
		r.Target == record.Target &&
		maps.Equal(r.Attributes, record.Attributes)
}
