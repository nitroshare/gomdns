package util

// MapsEqual compares two maps of identical types for equality.
func MapsEqual[K comparable, V comparable](a, b map[K]V) bool {
	if len(a) != len(b) {
		return false
	}
	for kA, vA := range a {
		vB, ok := b[kA]
		if !ok || vA != vB {
			return false
		}
	}
	return true
}
