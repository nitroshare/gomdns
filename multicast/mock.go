package multicast

func reset() {
	Interfaces = netInterfaces
	listenMulticastUDP = netListenMulticastUDP
}

func init() {
	reset()
}

// Mock replaces all internal functions with mocked equivalents.
func Mock() {
	Interfaces = mockInterfaces
	listenMulticastUDP = mockListenMulticastUDP
}

// MockWithError replaces all internal functions with equivalents that throw
// errors when invoked.
func MockWithError() {
	Interfaces = mockInterfacesWithError
	listenMulticastUDP = mockListenMulticastUDPWithError
}

// Unmock undoes the actions of Mock.
func Unmock() {
	reset()
}
