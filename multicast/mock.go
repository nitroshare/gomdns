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

// MockWithError replaces Interfaces with an equivalent that throws an error
// when invoked.
func MockWithError() {
	Interfaces = mockInterfacesWithError
	listenMulticastUDP = mockListenMulticastUDP
}

// Unmock undoes the actions of Mock.
func Unmock() {
	reset()
}
