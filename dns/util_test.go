package dns

import (
	"testing"

	"github.com/nitroshare/compare"
)

func TestParseName(t *testing.T) {
	for _, v := range []struct {
		Name        string
		Input       []byte
		Output      string
		StartOffset int
		EndOffset   int
		Err         bool
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
			Name:      "Invalid offset",
			Input:     []byte{1},
			EndOffset: 1,
			Err:       true,
		},
		{
			Name:        "Indirect",
			Input:       []byte{0, 1, 'y', 0, 1, 'x', 0xc0, 1},
			Output:      "x.y.",
			StartOffset: 4,
			EndOffset:   8,
		},
		{
			Name:      "Invalid indirect offset",
			Input:     []byte{0xc0},
			EndOffset: 1,
			Err:       true,
		},
		{
			Name:        "Infinite recursion",
			Input:       []byte{0, 1, 'x', 0xc0, 1},
			StartOffset: 1,
			EndOffset:   5,
			Err:         true,
		},
		{
			Name:      "Invalid first byte",
			Input:     []byte{0x80},
			EndOffset: 1,
			Err:       true,
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			offset := v.StartOffset
			n, err := parseName(v.Input, &offset)
			compare.Compare(t, n, v.Output, true)
			compare.Compare(t, offset, v.EndOffset, true)
			compare.Compare(t, err != nil, v.Err, true)
		})
	}
}
