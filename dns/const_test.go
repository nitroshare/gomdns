package dns

import "testing"

func TestTypeToString(t *testing.T) {
	for _, v := range []uint16{
		TypeA,
		TypePTR,
		TypeTXT,
		TypeSRV,
		TypeAAAA,
		TypeAny,
		254,
	} {
		TypeToString(v)
	}
}
