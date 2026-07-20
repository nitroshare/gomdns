package dns

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
)

// In order to aid in reducing the size of DNS packets, a technique called
// indirect pointers is used. Basically, a name can point to part of another
// name in a different part of the packet. Great care needs to be taken during
// parsing to make sure pointers are valid.

var (
	errSerializingName = errors.New("error serializing name")
	errParsingName     = errors.New("error parsing name")
)

// TODO: take advantage of indirect pointers to conserve bytes

func serializeName(name string) ([]byte, error) {
	var (
		b = bytes.Buffer{}
		v = []byte(name)
	)
	if len(v) > 0 {
		for _, v := range bytes.Split(v, []byte{'.'}) {
			if len(v) > 63 {
				return nil, errSerializingName
			}
			b.WriteByte(byte(len(v)))
			b.Write(v)
		}
	}
	b.WriteByte(0)
	return b.Bytes(), nil
}

func parseName(data []byte, offset *int) (string, error) {
	var (
		offsetEnd = 0
		offsetPtr = *offset
		name      strings.Builder
	)
	for {
		var numBytes uint8
		n, err := binary.Decode(data[*offset:], binary.BigEndian, &numBytes)
		if err != nil {
			return "", err
		}
		*offset += n
		if numBytes == 0 {
			break
		}
		switch numBytes & 0xc0 {
		case 0x00:
			n := int(numBytes)
			if *offset+n > len(data) {
				return "", errParsingName
			}
			name.WriteString(string(data[*offset : *offset+n]))
			name.WriteString(".")
			*offset += n
		case 0xc0:
			var numBytes2 uint8
			n, err := binary.Decode(data[*offset:], binary.BigEndian, &numBytes2)
			if err != nil {
				return "", err
			}
			*offset += n
			newOffset := int((uint16(numBytes & ^uint8(0xc0)) << 8) | uint16(numBytes2))
			if newOffset >= offsetPtr {
				return "", errParsingName
			}
			offsetPtr = newOffset
			if offsetEnd == 0 {
				offsetEnd = *offset
			}
			*offset = newOffset
		default:
			return "", errParsingName
		}
	}
	if offsetEnd != 0 {
		*offset = offsetEnd
	}
	return name.String(), nil
}
