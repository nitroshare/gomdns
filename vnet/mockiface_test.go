package vnet

import (
	"testing"
)

func TestMockInterface(t *testing.T) {
	i := NewMockInterface()
	i.Name()
	i.Addrs()
	i.Flags()
}
