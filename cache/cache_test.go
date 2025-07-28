package cache

import (
	"testing"
	"time"

	"github.com/nitroshare/mocktime"
)

func init() {
	chanTest = make(chan any)
}

func TestRecordQueryExpiry(t *testing.T) {
	mocktime.Mock()
	defer mocktime.Unmock()
	var (
		chanQuery   = make(chan *Record)
		chanExpired = make(chan *Record)
		c           = New(&Config{
			ChanQuery:   chanQuery,
			ChanExpired: chanExpired,
		})
	)
	r := &Record{
		TTL: 10,
	}
	c.Add(r)
	mocktime.Advance(6 * time.Second)
	<-chanTest
	<-chanQuery
	mocktime.Advance(5 * time.Second)
	<-chanTest
	<-chanExpired
}
