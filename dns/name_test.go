package dns

import (
	"reflect"
	"strings"
	"testing"

	"github.com/nitroshare/compare"
)

func TestSerializeName(t *testing.T) {
	for _, v := range []struct {
		Name   string
		Input  string
		Output []byte
		Err    bool
	}{
		{
			Name:   "Empty name",
			Output: []byte{0},
		},
		{
			Name:  "Valid name",
			Input: "x",
			Output: []byte{
				1, 'x', 0,
			},
		},
		{
			Name:  "Name > 63 bytes",
			Input: strings.Repeat("0", 64),
			Err:   true,
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			n, err := serializeName(v.Input)
			compare.Compare(t, reflect.DeepEqual(n, v.Output), true, true)
			compare.Compare(t, err != nil, v.Err, true)
		})
	}
}

func TestParseName(t *testing.T) {
	for _, v := range []struct {
		Name        string
		Input       []byte
		StartOffset int
		Output      string
		Err         bool
		EndOffset   int
	}{
		{
			Name:  "Empty data",
			Input: []byte{},
			Err:   true,
		},
		{
			Name:      "Simple name",
			Input:     []byte{1, 'x', 1, 'x', 0},
			Output:    "x.x.",
			EndOffset: 5,
		},
		{
			Name:  "Invalid offset",
			Input: []byte{1},
			Err:   true,
		},
		{
			Name:        "Indirect",
			Input:       []byte{0, 1, 'y', 0, 1, 'x', 0xc0, 1},
			StartOffset: 4,
			Output:      "x.y.",
			EndOffset:   8,
		},
		{
			Name:  "Invalid indirect offset",
			Input: []byte{0xc0},
			Err:   true,
		},
		{
			Name:        "Infinite recursion",
			Input:       []byte{0, 1, 'x', 0xc0, 1},
			StartOffset: 1,
			Err:         true,
		},
		{
			Name:  "Invalid first byte",
			Input: []byte{0x80},
			Err:   true,
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			offset := v.StartOffset
			n, err := parseName(v.Input, &offset)
			compare.Compare(t, n, v.Output, true)
			compare.Compare(t, err != nil, v.Err, true)
			if v.EndOffset != 0 {
				compare.Compare(t, offset, v.EndOffset, true)
			}
		})
	}
}
