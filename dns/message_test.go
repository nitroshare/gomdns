package dns

import (
	"net/netip"
	"reflect"
	"strings"
	"testing"

	"github.com/nitroshare/compare"
)

func TestMessageString(t *testing.T) {
	compare.Compare(
		t,
		(&Message{
			TransactionID: 1,
		}).String(),
		"query id:1",
		true,
	)
	compare.Compare(
		t,
		(&Message{
			Response: true,
		}).String(),
		"response id:0",
		true,
	)
}

func TestMessageSerialize(t *testing.T) {
	for _, v := range []struct {
		Name   string
		Input  *Message
		Output []byte
		Err    bool
	}{
		{
			Name: "Valid message",
			Input: &Message{
				TransactionID: 1,
				Response:      true,
				Truncated:     true,
				Questions: []*Question{
					{
						Name: "x",
						Type: TypeA,
					},
				},
				Records: []*Record{
					{
						Name:    "x",
						Type:    TypeA,
						Address: netip.MustParseAddr("255.255.0.0"),
					},
				},
			},
			Output: []byte{
				0, 1,
				0x86, 0,
				0, 1,
				0, 1,
				0, 0,
				0, 0,
				1, 'x', 0,
				0, 1,
				0, 0,
				1, 'x', 0,
				0, 1,
				0, 0,
				0, 0, 0, 0,
				0, 4,
				255, 255, 0, 0,
			},
		},
		{
			Name: "Invalid question",
			Input: &Message{
				Questions: []*Question{
					{Name: strings.Repeat("0", 64)},
				},
			},
			Err: true,
		},
		{
			Name: "Invalid record",
			Input: &Message{
				Records: []*Record{
					{Name: strings.Repeat("0", 64)},
				},
			},
			Err: true,
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			b, err := v.Input.Serialize()
			compare.Compare(t, reflect.DeepEqual(b, v.Output), true, true)
			compare.Compare(t, err != nil, v.Err, true)
		})
	}
}

func TestParseMessage(t *testing.T) {
	for _, v := range []struct {
		Name   string
		Input  []byte
		Output *Message
		Err    bool
	}{
		{
			Name: "Empty message",
			Err:  true,
		},
		{
			Name: "Valid message",
			Input: []byte{
				0, 0,
				0x86, 0,
				0, 1,
				0, 1,
				0, 0,
				0, 0,
				1, 'x', 0, 0, 1, 0, 0,
				0xc0, 0x0c,
				0, 1,
				0, 0,
				0, 0, 0, 0,
				0, 4,
				0, 0, 0, 0,
			},
			Output: &Message{
				Response:  true,
				Truncated: true,
				Questions: []*Question{
					{
						Name: "x.",
						Type: TypeA,
					},
				},
				Records: []*Record{
					{
						Name:    "x.",
						Type:    TypeA,
						Address: netip.MustParseAddr("0.0.0.0"),
					},
				},
			},
		},
		{
			Name: "Invalid question",
			Input: []byte{
				0, 0,
				0x86, 0,
				0, 1,
				0, 0,
				0, 0,
				0, 0,
			},
			Err: true,
		},
		{
			Name: "Invalid record",
			Input: []byte{
				0, 0,
				0x86, 0,
				0, 0,
				0, 1,
				0, 0,
				0, 0,
			},
			Err: true,
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			m, err := ParseMessage(v.Input)
			compare.Compare(t, reflect.DeepEqual(m, v.Output), true, true)
			compare.Compare(t, err != nil, v.Err, true)
		})
	}
}
