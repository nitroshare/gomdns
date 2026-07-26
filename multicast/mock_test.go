package multicast

import (
	"testing"

	"github.com/nitroshare/gomdns/compare"
)

func TestMockUnmock(t *testing.T) {
	compare.CompareFn(t, Interfaces, netInterfaces, true)
	compare.CompareFn(t, listenMulticastUDP, netListenMulticastUDP, true)
	Mock()
	compare.CompareFn(t, Interfaces, mockInterfaces, true)
	compare.CompareFn(t, listenMulticastUDP, mockListenMulticastUDP, true)
	MockWithError()
	compare.CompareFn(t, Interfaces, mockInterfacesWithError, true)
	compare.CompareFn(t, listenMulticastUDP, mockListenMulticastUDPWithError, true)
	Unmock()
	compare.CompareFn(t, Interfaces, netInterfaces, true)
	compare.CompareFn(t, listenMulticastUDP, netListenMulticastUDP, true)
}
