package dns

import (
	"encoding/binary"
	"fmt"
)

type questionFields struct {
	Type  uint16
	Class uint16
}

// Question represents a DNS query for records.
type Question struct {
	Name    string
	Type    uint16
	Unicast bool
}

func (q Question) String() string {
	return fmt.Sprintf("%s %s", TypeToString(q.Type), q.Name)
}

func parseQuestion(data []byte, offset *int) (*Question, error) {
	v, err := parseName(data, offset)
	if err != nil {
		return nil, err
	}
	var fields questionFields
	n, err := binary.Decode(data[*offset:], binary.BigEndian, &fields)
	if err != nil {
		return nil, err
	}
	*offset += n
	return &Question{
		Name:    v,
		Type:    fields.Type,
		Unicast: fields.Class&0x8000 != 0,
	}, nil
}
