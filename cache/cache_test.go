package cache

import (
	"reflect"
	"testing"

	"github.com/nitroshare/compare"
	"github.com/nitroshare/gomdns/dns"
	"github.com/nitroshare/gotime"
)

const (
	testName = "name"
	testType = dns.TypeA
	testTTL  = 10
)

var (
	testRecord = &dns.Record{
		Name: testName,
		Type: testType,
		Ttl:  testTTL,
	}
)

func TestQueryAndExpiry(t *testing.T) {
	gotime.Mock()
	defer gotime.Unmock()
	var (
		chanQuery   = make(chan *dns.Record)
		chanExpired = make(chan *dns.Record)
		c           = New(&Config{
			ChanQuery:   chanQuery,
			ChanExpired: chanExpired,
		})
	)
	defer c.Close()
	c.Add(testRecord)
	for range 4 {
		gotime.AdvanceToAfter()
		<-chanQuery
	}
	gotime.AdvanceToAfter()
	<-chanExpired
}

func TestLookup(t *testing.T) {
	gotime.Mock()
	defer gotime.Unmock()
	c := New(&Config{})
	defer c.Close()
	for range 2 {
		c.Add(testRecord)
		compare.Compare(
			t,
			reflect.DeepEqual(
				c.Lookup(testName, testType),
				[]*dns.Record{testRecord},
			),
			true,
			true,
		)
	}
}

func TestFlush(t *testing.T) {
	gotime.Mock()
	defer gotime.Unmock()
	var (
		chanExpired = make(chan *dns.Record)
		c           = New(&Config{
			ChanExpired: chanExpired,
		})
	)
	defer c.Close()
	c.Add(testRecord)
	go c.Add(&dns.Record{
		Name:       testName,
		Type:       testType,
		FlushCache: true,
	})
	<-chanExpired
	compare.Compare(
		t,
		reflect.DeepEqual(
			c.Lookup(testName, testType),
			[]*dns.Record{},
		),
		true,
		true,
	)
}

func TestNonBlockingSend(t *testing.T) {
	gotime.Mock()
	defer gotime.Unmock()
	var (
		chanQuery   = make(chan *dns.Record)
		chanExpired = make(chan *dns.Record)
		c           = New(&Config{
			ChanQuery:   chanQuery,
			ChanExpired: chanExpired,
		})
	)
	defer func() { <-chanQuery }()
	defer c.Close()
	c.Add(testRecord)
	for range 5 {
		gotime.AdvanceToAfter()
	}
	<-chanExpired
}
