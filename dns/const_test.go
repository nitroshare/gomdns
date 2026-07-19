package dns

import "testing"

func TestTypeToString(t *testing.T) {
	for _, v := range []int{
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
