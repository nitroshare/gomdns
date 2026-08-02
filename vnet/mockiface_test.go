package vnet

import (
	"errors"
	"testing"

	"github.com/nitroshare/gomdns/compare"
)

func TestMockInterface(t *testing.T) {
	i := NewMockInterface()
	i.Name()
	i.Addrs()
	i.Flags()
	i.MockListenError = errors.New("test")
	_, e := i.Listen("", nil)
	compare.Compare(t, e != nil, true, true)
}
