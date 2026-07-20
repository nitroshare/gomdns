package dns

import (
	"net/netip"
	"reflect"
	"strings"
	"testing"

	"github.com/nitroshare/compare"
)

func TestRecordString(t *testing.T) {
	compare.Compare(
		t,
		(&Record{
			Name:    "x.",
			Type:    TypeA,
			Address: netip.MustParseAddr("255.255.0.0"),
		}).String(),
		"A x. 255.255.0.0",
		true,
	)
	compare.Compare(
		t,
		(&Record{
			Name:   "x.",
			Type:   TypePTR,
			Target: "y.",
		}).String(),
		"PTR x. y.",
		true,
	)
	compare.Compare(
		t,
		(&Record{
			Name: "x.",
			Type: TypeTXT,
			Attributes: []string{
				"1",
				"2",
			},
		}).String(),
		"TXT x. 1, 2",
		true,
	)
	compare.Compare(
		t,
		(&Record{
			Name:     "x.",
			Type:     TypeSRV,
			Priority: 1,
			Weight:   2,
			Port:     80,
			Target:   "y.",
		}).String(),
		"SRV x. priority=1 weight=2 port=80 y.",
		true,
	)
	compare.Compare(
		t,
		(&Record{
			Name: "x.",
		}).String(),
		"?? x.",
		true,
	)
}

func TestRecordSerialize(t *testing.T) {
	for _, v := range []struct {
		Name   string
		Input  *Record
		Output []byte
		Err    bool
	}{
		{
			Name: "Invalid record name",
			Input: &Record{
				Name: strings.Repeat("0", 64),
			},
			Err: true,
		},
		{
			Name: "Valid A record",
			Input: &Record{
				Name:       "x",
				Type:       TypeA,
				FlushCache: true,
				Ttl:        1,
				Address:    netip.MustParseAddr("255.255.0.0"),
			},
			Output: []byte{
				1, 'x', 0,
				0, 1,
				0x80, 0,
				0, 0, 0, 1,
				0, 4,
				255, 255, 0, 0,
			},
		},
		{
			Name: "Valid PTR record",
			Input: &Record{
				Name:   "x",
				Type:   TypePTR,
				Target: "y",
			},
			Output: []byte{
				1, 'x', 0,
				0, 12,
				0, 0,
				0, 0, 0, 0,
				0, 3,
				1, 'y', 0,
			},
		},
		{
			Name: "Valid TXT record",
			Input: &Record{
				Name:       "x",
				Type:       TypeTXT,
				Attributes: []string{"y"},
			},
			Output: []byte{
				1, 'x', 0,
				0, 16,
				0, 0,
				0, 0, 0, 0,
				0, 2,
				1, 'y',
			},
		},
		{
			Name: "Valid SRV record",
			Input: &Record{
				Name:     "x",
				Type:     TypeSRV,
				Priority: 1,
				Weight:   2,
				Port:     80,
				Target:   "y",
			},
			Output: []byte{
				1, 'x', 0,
				0, 33,
				0, 0,
				0, 0, 0, 0,
				0, 9,
				0, 1,
				0, 2,
				0, 80,
				1, 'y', 0,
			},
		},
		{
			Name: "Valid AAAA record",
			Input: &Record{
				Name:    "x",
				Type:    TypeAAAA,
				Address: netip.MustParseAddr("::1"),
			},
			Output: []byte{
				1, 'x', 0,
				0, 28,
				0, 0,
				0, 0, 0, 0,
				0, 16,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
			},
		},
		{
			Name: "Invalid PTR record",
			Input: &Record{
				Type:   TypePTR,
				Target: strings.Repeat("0", 64),
			},
			Err: true,
		},
		{
			Name: "Invalid TXT record",
			Input: &Record{
				Type:       TypeTXT,
				Attributes: []string{strings.Repeat("0", 64)},
			},
			Err: true,
		},
		{
			Name: "Invalid SRV record",
			Input: &Record{
				Type:   TypeSRV,
				Target: strings.Repeat("0", 64),
			},
			Err: true,
		},
		{
			Name: "Unknown record type",
			Input: &Record{
				Type: 255,
			},
			Err: true,
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			b, err := v.Input.serialize()
			compare.Compare(t, reflect.DeepEqual(b, v.Output), true, true)
			compare.Compare(t, err != nil, v.Err, true)
		})
	}
}

func TestParseRecord(t *testing.T) {
	for _, v := range []struct {
		Name        string
		Input       []byte
		StartOffset int
		Output      *Record
		Err         bool
		EndOffset   int
	}{
		{
			Name: "Empty record",
			Err:  true,
		},
		{
			Name: "Valid A record",
			Input: []byte{
				1, 'x', 0,
				0, 1,
				0x80, 0,
				0, 0, 0, 60,
				0, 4,
				0xff, 0xff, 0, 0,
			},
			Output: &Record{
				Name:       "x.",
				Type:       TypeA,
				FlushCache: true,
				Ttl:        60,
				Address:    netip.MustParseAddr("255.255.0.0"),
			},
			EndOffset: 17,
		},
		{
			Name: "Valid PTR record",
			Input: []byte{
				1, 'x', 0,
				0, 12,
				0, 0,
				0, 0, 0, 0,
				0, 3,
				1, 'y', 0,
			},
			Output: &Record{
				Name:   "x.",
				Type:   TypePTR,
				Target: "y.",
			},
			EndOffset: 16,
		},
		{
			Name: "Valid TXT record",
			Input: []byte{
				1, 'x', 0,
				0, 16,
				0, 0,
				0, 0, 0, 0,
				0, 5,
				1, 'x', 1, 'y', 0,
			},
			Output: &Record{
				Name: "x.",
				Type: TypeTXT,
				Attributes: []string{
					"x",
					"y",
					"",
				},
			},
			EndOffset: 18,
		},
		{
			Name: "Valid SRV record",
			Input: []byte{
				1, 'x', 0,
				0, 33,
				0, 0,
				0, 0, 0, 0,
				0, 9,
				0, 1,
				0, 2,
				0, 80,
				1, 'y', 0,
			},
			Output: &Record{
				Name:     "x.",
				Type:     TypeSRV,
				Priority: 1,
				Weight:   2,
				Port:     80,
				Target:   "y.",
			},
			EndOffset: 22,
		},
		{
			Name: "Valid AAAA record",
			Input: []byte{
				1, 'x', 0,
				0, 28,
				0, 0,
				0, 0, 0, 0,
				0, 16,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
			},
			Output: &Record{
				Name:    "x.",
				Type:    TypeAAAA,
				Address: netip.MustParseAddr("::1"),
			},
			EndOffset: 29,
		},
		{
			Name: "Truncated (after name)",
			Input: []byte{
				1, 'x', 0,
			},
			Err: true,
		},
		{
			Name: "Invalid A record",
			Input: []byte{
				1, 'x', 0,
				0, 1,
				0, 0,
				0, 0, 0, 0,
				0, 0,
			},
			Err: true,
		},
		{
			Name: "Invalid PTR record",
			Input: []byte{
				1, 'x', 0,
				0, 12,
				0, 0,
				0, 0, 0, 0,
				0, 0,
			},
			Err: true,
		},
		{
			Name: "Invalid TXT record (missing byte len)",
			Input: []byte{
				1, 'x', 0,
				0, 16,
				0, 0,
				0, 0, 0, 0,
				0, 1,
			},
			Err: true,
		},
		{
			Name: "Invalid TXT record (invalid byte len)",
			Input: []byte{
				1, 'x', 0,
				0, 16,
				0, 0,
				0, 0, 0, 0,
				0, 2,
				2, 0,
			},
			Err: true,
		},
		{
			Name: "Invalid SRV record (truncated)",
			Input: []byte{
				1, 'x', 0,
				0, 33,
				0, 0,
				0, 0, 0, 0,
				0, 0,
			},
			Err: true,
		},
		{
			Name: "Invalid SRV record (invalid name)",
			Input: []byte{
				1, 'x', 0,
				0, 33,
				0, 0,
				0, 0, 0, 0,
				0, 6,
				0, 0,
				0, 0,
				0, 0,
			},
			Err: true,
		},
		{
			Name: "Invalid AAAA record",
			Input: []byte{
				1, 'x', 0,
				0, 28,
				0, 0,
				0, 0, 0, 0,
				0, 0,
			},
			Err: true,
		},
		{
			Name: "Offset beyond packet",
			Input: []byte{
				1, 'x', 0,
				0, 255,
				0, 0,
				0, 0, 0, 0,
				0, 1,
			},
			Err: true,
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			offset := v.StartOffset
			r, err := parseRecord(v.Input, &offset)
			compare.Compare(t, reflect.DeepEqual(r, v.Output), true, true)
			compare.Compare(t, err != nil, v.Err, true)
			if v.EndOffset != 0 {
				compare.Compare(t, offset, v.EndOffset, true)
			}
		})
	}
}
