package cache

import (
	"reflect"
	"testing"

	"github.com/miekg/dns"
	"github.com/nitroshare/compare"
	"github.com/nitroshare/mocktime"
)

const (
	testName = "name"
	testType = dns.TypeA
	testTTL  = 10
)

var (
	testRecord = &Record{
		Name: testName,
		Type: testType,
		TTL:  testTTL,
	}
)

func TestRecordQueryAndExpiry(t *testing.T) {
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
	defer c.Close()
	c.Add(testRecord)
	for range 4 {
		mocktime.AdvanceToAfter()
		<-chanQuery
	}
	mocktime.AdvanceToAfter()
	<-chanExpired
}

func TestRecordLookup(t *testing.T) {
	mocktime.Mock()
	defer mocktime.Unmock()
	c := New(&Config{})
	defer c.Close()
	for range 2 {
		c.Add(testRecord)
		compare.Compare(
			t,
			reflect.DeepEqual(
				c.Lookup(testName, testType),
				[]*Record{testRecord},
			),
			true,
			true,
		)
	}
}

func TestRecordFlush(t *testing.T) {
	mocktime.Mock()
	defer mocktime.Unmock()
	var (
		chanExpired = make(chan *Record)
		c           = New(&Config{
			ChanExpired: chanExpired,
		})
	)
	defer c.Close()
	c.Add(testRecord)
	go c.Add(&Record{
		Name:       testName,
		Type:       testType,
		FlushCache: true,
	})
	<-chanExpired
	compare.Compare(
		t,
		reflect.DeepEqual(
			c.Lookup(testName, testType),
			[]*Record{},
		),
		true,
		true,
	)
}
