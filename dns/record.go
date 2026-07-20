package dns

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net/netip"
	"strings"
)

var (
	errSerializingRecord = errors.New("error serializing record")
	errParsingRecord     = errors.New("error parsing record")
)

type recordFields struct {
	Type    uint16
	Class   uint16
	Ttl     uint32
	DataLen uint16
}

type recordSRVFields struct {
	Priority uint16
	Weight   uint16
	Port     uint16
}

// Record represents an answer to a DNS query. Not all fields are used by each
// type of record, though the first four fields are common to all records.
type Record struct {
	Name       string
	Type       uint16
	FlushCache bool
	Ttl        uint32
	Address    netip.Addr
	Target     string
	Attributes []string
	Priority   uint16
	Weight     uint16
	Port       uint16
}

func (r Record) String() string {
	v := fmt.Sprintf("%s %s", TypeToString(r.Type), r.Name)
	switch r.Type {
	case TypeA, TypeAAAA:
		return fmt.Sprintf("%s %s", v, r.Address)
	case TypePTR:
		return fmt.Sprintf("%s %s", v, r.Target)
	case TypeTXT:
		return fmt.Sprintf("%s %s", v, strings.Join(r.Attributes, ", "))
	case TypeSRV:
		return fmt.Sprintf(
			"%s priority=%d weight=%d port=%d %s",
			v,
			r.Priority,
			r.Weight,
			r.Port,
			r.Target,
		)
	default:
		return v
	}
}

func (r Record) serialize() ([]byte, error) {
	b := &bytes.Buffer{}
	n, err := serializeName(r.Name)
	if err != nil {
		return nil, err
	}
	b.Write(n)
	fields := &recordFields{
		Type: r.Type,
		Ttl:  r.Ttl,
	}
	if r.FlushCache {
		fields.Class |= 0x8000
	}
	binary.Write(b, binary.BigEndian, fields)
	offset := b.Len()
	switch r.Type {
	case TypeA:
		binary.Write(b, binary.BigEndian, r.Address.As4())
	case TypePTR:
		n, err := serializeName(r.Target)
		if err != nil {
			return nil, err
		}
		b.Write(n)
	case TypeTXT:
		for _, a := range r.Attributes {
			v := []byte(a)
			if len(v) > 63 {
				return nil, errSerializingName
			}
			b.WriteByte(uint8(len(v)))
			b.Write(v)
		}
	case TypeSRV:
		fields := &recordSRVFields{
			Priority: r.Priority,
			Weight:   r.Weight,
			Port:     r.Port,
		}
		binary.Write(b, binary.BigEndian, fields)
		n, err := serializeName(r.Target)
		if err != nil {
			return nil, err
		}
		b.Write(n)
	case TypeAAAA:
		binary.Write(b, binary.BigEndian, r.Address.As16())
	default:
		return nil, errSerializingRecord
	}
	v := b.Bytes()
	binary.BigEndian.PutUint16(v[offset-2:], uint16(b.Len()-offset))
	return v, nil
}

func parseRecord(data []byte, offset *int) (*Record, error) {
	v, err := parseName(data, offset)
	if err != nil {
		return nil, err
	}
	var fields recordFields
	n, err := binary.Decode(data[*offset:], binary.BigEndian, &fields)
	if err != nil {
		return nil, err
	}
	*offset += n
	r := &Record{
		Name:       v,
		Type:       fields.Type,
		FlushCache: fields.Class&0x8000 != 0,
		Ttl:        fields.Ttl,
	}
	offsetStart := *offset
	switch fields.Type {
	case TypeA:
		var v [4]byte
		n, err := binary.Decode(data[*offset:], binary.BigEndian, &v)
		if err != nil {
			return nil, err
		}
		*offset += n
		r.Address = netip.AddrFrom4(v)
	case TypePTR:
		v, err := parseName(data, offset)
		if err != nil {
			return nil, err
		}
		r.Target = v
	case TypeTXT:
		offsetPtr := *offset
		for *offset < offsetPtr+int(fields.DataLen) {
			var numBytes uint8
			n, err := binary.Decode(data[*offset:], binary.BigEndian, &numBytes)
			if err != nil {
				return nil, err
			}
			*offset += n
			v := int(numBytes)
			if v == 0 {
				r.Attributes = append(r.Attributes, "")
				continue
			}
			if *offset+v > len(data) {
				return nil, errParsingRecord
			}
			r.Attributes = append(r.Attributes, string(data[*offset:*offset+v]))
			*offset += v
		}
	case TypeSRV:
		var fields recordSRVFields
		n, err := binary.Decode(data[*offset:], binary.BigEndian, &fields)
		if err != nil {
			return nil, err
		}
		*offset += n
		r.Priority = fields.Priority
		r.Weight = fields.Weight
		r.Port = fields.Port
		v, err := parseName(data, offset)
		if err != nil {
			return nil, err
		}
		r.Target = v
	case TypeAAAA:
		var v [16]byte
		n, err := binary.Decode(data[*offset:], binary.BigEndian, &v)
		if err != nil {
			return nil, err
		}
		*offset += n
		r.Address = netip.AddrFrom16(v)
	}
	*offset = offsetStart + int(fields.DataLen)
	if *offset > len(data) {
		return nil, errParsingRecord
	}
	return r, nil
}
