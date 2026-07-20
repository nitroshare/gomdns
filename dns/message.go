package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net/netip"
)

type messageFields struct {
	TransactionID uint16
	Flags         uint16
	NumQuestions  uint16
	NumAnswers    uint16
	NumAuthority  uint16
	NumAdditional uint16
}

// Message represents a DNS message to be sent or received over the wire.
type Message struct {
	Address       netip.Addr
	Port          uint16
	TransactionID uint16
	Response      bool
	Truncated     bool
	Questions     []*Question
	Records       []*Record
}

func (m Message) String() string {
	var type_ string
	if m.Response {
		type_ = "response"
	} else {
		type_ = "query"
	}
	return fmt.Sprintf(
		"%s id:%d",
		type_,
		m.TransactionID,
	)
}

// Serialize converts the message into its wire format.
func (m Message) Serialize() ([]byte, error) {
	b := &bytes.Buffer{}
	fields := &messageFields{
		TransactionID: m.TransactionID,
		NumQuestions:  uint16(len(m.Questions)),
		NumAnswers:    uint16(len(m.Records)),
	}
	if m.Response {
		fields.Flags |= 0x8400
	}
	if m.Truncated {
		fields.Flags |= 0x0200
	}
	binary.Write(b, binary.BigEndian, fields)
	for _, q := range m.Questions {
		v, err := q.serialize()
		if err != nil {
			return nil, err
		}
		b.Write(v)
	}
	for _, r := range m.Records {
		v, err := r.serialize()
		if err != nil {
			return nil, err
		}
		b.Write(v)
	}
	return b.Bytes(), nil
}

// ParseMessage attempts to parse a DNS message. Note that this does not know
// the origin, and therefore cannot fill in the Address and Port fields.
func ParseMessage(data []byte) (*Message, error) {
	var fields messageFields
	n, err := binary.Decode(data, binary.BigEndian, &fields)
	if err != nil {
		return nil, err
	}
	var (
		msg = &Message{
			TransactionID: fields.TransactionID,
			Response:      fields.Flags&0x8400 != 0,
			Truncated:     fields.Flags&0x0200 != 0,
		}
		offset = n
	)
	for i := uint16(0); i < fields.NumQuestions; i++ {
		q, err := parseQuestion(data, &offset)
		if err != nil {
			return nil, err
		}
		msg.Questions = append(msg.Questions, q)
	}
	for range fields.NumAnswers + fields.NumAuthority + fields.NumAdditional {
		r, err := parseRecord(data, &offset)
		if err != nil {
			return nil, err
		}
		msg.Records = append(msg.Records, r)
	}
	return msg, nil
}
