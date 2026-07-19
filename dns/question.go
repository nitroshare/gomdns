package dns

import (
	"encoding/binary"
	"fmt"
)

type questionFields struct {
	Type  uint16
	Class uint16
}

type Question struct {
	Name  string
	Type  int
	Class int
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
		Name:  v,
		Type:  int(fields.Type),
		Class: int(fields.Class),
	}, nil
}
