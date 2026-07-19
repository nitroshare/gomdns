package dns

import (
	"reflect"
	"testing"

	"github.com/nitroshare/compare"
)

func TestQuestionToString(t *testing.T) {
	compare.Compare(
		t,
		Question{
			Name: "x.",
			Type: TypeA,
		}.String(),
		"A x.",
		true,
	)
}

func TestParseQuestion(t *testing.T) {
	for _, v := range []struct {
		Name        string
		Input       []byte
		Output      *Question
		StartOffset int
		EndOffset   int
		Err         bool
	}{
		{
			Name:  "Empty input",
			Input: []byte{},
			Err:   true,
		},
		{
			Name:  "Valid question",
			Input: []byte{1, 'x', 0, 0, 1, 0x80, 0},
			Output: &Question{
				Name:    "x.",
				Type:    TypeA,
				Unicast: true,
			},
			EndOffset: 7,
		},
		{
			Name:      "Truncated question",
			Input:     []byte{1, 'x', 0},
			EndOffset: 3,
			Err:       true,
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			offset := v.StartOffset
			q, err := parseQuestion(v.Input, &offset)
			compare.Compare(t, reflect.DeepEqual(q, v.Output), true, true)
			compare.Compare(t, offset, v.EndOffset, true)
			compare.Compare(t, err != nil, v.Err, true)
		})
	}
}
