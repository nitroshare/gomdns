package dns

const (
	TypeA    = 1
	TypePTR  = 12
	TypeTXT  = 16
	TypeSRV  = 33
	TypeAAAA = 28
	TypeAny  = 255
)

// TypeToString returns the string value for the specified type.
func TypeToString(type_ int) string {
	switch type_ {
	case TypeA:
		return "A"
	case TypePTR:
		return "PTR"
	case TypeTXT:
		return "TXT"
	case TypeSRV:
		return "SRV"
	case TypeAAAA:
		return "AAAA"
	case TypeAny:
		return "*"
	default:
		return "??"
	}
}
