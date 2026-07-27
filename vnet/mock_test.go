package vnet

import (
	"testing"

	"github.com/nitroshare/gomdns/compare"
)

func TestMockUnmock(t *testing.T) {
	compare.CompareFn(t, Interfaces, netInterfaces, true)
	Mock()
	compare.CompareFn(t, Interfaces, mockInterfaces, true)
	Unmock()
	compare.CompareFn(t, Interfaces, netInterfaces, true)
}

func TestAddInterface(t *testing.T) {
	Mock()
	defer Unmock()
	AddInterface(&netInterface{})
	i, err := Interfaces()
	if err != nil {
		t.Fatal(err)
	}
	compare.Compare(t, len(i), 1, true)
}
