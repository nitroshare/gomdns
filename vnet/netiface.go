package vnet

import (
	"net"
)

type netInterface struct {
	i *net.Interface
}

func (n netInterface) Addrs() ([]net.Addr, error) { return n.i.Addrs() }
func (n netInterface) Flags() net.Flags           { return n.i.Flags }
func (n netInterface) Interface() *net.Interface  { return n.i }
