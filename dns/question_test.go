package dns

import (
	"reflect"
	"strings"
	"testing"

	"github.com/nitroshare/compare"
)

func TestQuestionString(t *testing.T) {
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

func TestQuestionSerialize(t *testing.T) {
	for _, v := range []struct {
		Name   string
		Input  *Question
		Output []byte
		Err    bool
	}{
		{
			Name: "Invalid question name",
			Input: &Question{
				Name: strings.Repeat("0", 64),
			},
			Err: true,
		},
		{
			Name: "Valid question",
			Input: &Question{
				Name:    "x",
				Type:    TypeA,
				Unicast: true,
			},
			Output: []byte{
				1, 'x', 0,
				0, 1,
				0x80, 0,
			},
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			b, err := v.Input.serialize()
			compare.Compare(t, reflect.DeepEqual(b, v.Output), true, true)
			compare.Compare(t, err != nil, v.Err, true)
		})
	}
}

func TestParseQuestion(t *testing.T) {
	for _, v := range []struct {
		Name      string
		Input     []byte
		Output    *Question
		Err       bool
		EndOffset int
	}{
		{
			Name: "Empty question",
			Err:  true,
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
			Name:  "Truncated question",
			Input: []byte{1, 'x', 0},
			Err:   true,
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			offset := 0
			q, err := parseQuestion(v.Input, &offset)
			compare.Compare(t, reflect.DeepEqual(q, v.Output), true, true)
			compare.Compare(t, err != nil, v.Err, true)
			if v.EndOffset != 0 {
				compare.Compare(t, offset, v.EndOffset, true)
			}
		})
	}
}
