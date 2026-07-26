package cache

import (
	"reflect"
	"testing"

	"github.com/nitroshare/gomdns/compare"
	"github.com/nitroshare/gomdns/dns"
	"github.com/nitroshare/gomdns/vtime"
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
	vtime.Mock()
	defer vtime.Unmock()
	c := New(&Config{})
	defer c.Close()
	var (
		chanQuery   = c.Query.Subscribe()
		chanExpired = c.Expired.Subscribe()
	)
	c.Add(testRecord)
	for range 4 {
		vtime.AdvanceToAfter()
		<-chanQuery
	}
	vtime.AdvanceToAfter()
	<-chanExpired
}

func TestLookup(t *testing.T) {
	vtime.Mock()
	defer vtime.Unmock()
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
	vtime.Mock()
	defer vtime.Unmock()
	c := New(&Config{})
	defer c.Close()
	chanExpired := c.Expired.Subscribe()
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
