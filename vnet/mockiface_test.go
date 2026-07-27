package vnet

import (
	"testing"
)

func TestMockInterface(t *testing.T) {
	i := NewMockInterface()
	i.Addrs()
	i.Flags()

	//...
}
