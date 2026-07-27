package vnet

var (

	// Interfaces returns a list of network interfaces.
	Interfaces func() ([]Interface, error)

	interfaceList []Interface
)

func reset() {
	Interfaces = netInterfaces
}

func init() {
	reset()
}

// Mock replaces all public net functions with mocked equivalents.
func Mock() {
	Interfaces = mockInterfaces
	ClearInterfaces()
}

// Unmock undoes the actions of Mock.
func Unmock() {
	reset()
}

// AddInterfaces adds an interface to the list returned by Interfaces when
// Mock has been called.
func AddInterface(i Interface) {
	interfaceList = append(interfaceList, i)
}

// ClearInterfaces clears the list returned by Interfaces when Mock has been
// called.
func ClearInterfaces() {
	interfaceList = []Interface{}
}
