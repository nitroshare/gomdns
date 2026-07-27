package vnet

import (
	"net"
)

var (
	origNetInterfaces = net.Interfaces
)

func netInterfaces() ([]Interface, error) {
	l, err := origNetInterfaces()
	if err != nil {
		return nil, err
	}
	interfaces := []Interface{}
	for _, i := range l {
		interfaces = append(interfaces, &netInterface{&i})
	}
	return interfaces, nil
}

func mockInterfaces() ([]Interface, error) {
	return interfaceList, nil
}
