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

// Unmock undoes the actions of Mock.
func Unmock() {
	reset()
}
