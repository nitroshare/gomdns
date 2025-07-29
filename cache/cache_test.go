package cache

import (
	"testing"

	"github.com/nitroshare/mocktime"
)

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
	c.Add(&Record{
		TTL: 1,
	})
	for range 4 {
		mocktime.AdvanceToAfter()
		<-chanQuery
	}
	mocktime.AdvanceToAfter()
	<-chanExpired
}
