package compare

import (
	"testing"
)

type mockT struct {
	fatal bool
}

func (m *mockT) Helper()               {}
func (m *mockT) Fatalf(string, ...any) { m.fatal = true }

func TestCompare(t *testing.T) {
	for _, v := range []struct {
		name    string
		v2      bool
		matches bool
		fatal   bool
	}{
		{
			name:    "ShouldNotMatchAndDoesnt",
			v2:      false,
			matches: false,
			fatal:   false,
		},
		{
			name:    "ShouldNotMatchAndDoes",
			v2:      true,
			matches: false,
			fatal:   true,
		},
		{
			name:    "ShouldMatchAndDoesnt",
			v2:      false,
			matches: true,
			fatal:   true,
		},
		{
			name:    "ShouldMatchAndDoes",
			v2:      true,
			matches: true,
			fatal:   false,
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			m := &mockT{}
			Compare(m, true, v.v2, v.matches)
			if m.fatal != v.fatal {
				t.Fatalf("%v != %v", m.fatal, v.fatal)
			}
		})
	}
}

var (
	fn1   = func() int { return 1 }
	fn2   = func() int { return 2 }
	fnBad = 42
)

func TestCompareFn(t *testing.T) {
	for _, v := range []struct {
		name    string
		v1      any
		v2      any
		matches bool
		fatal   bool
	}{
		{
			name:    "identical",
			v1:      fn1,
			v2:      fn1,
			matches: true,
		},
		{
			name:    "not identical",
			v1:      fn1,
			v2:      fn2,
			matches: false,
		},
		{
			name:  "not function",
			v1:    fn1,
			v2:    fnBad,
			fatal: true,
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			m := &mockT{}
			CompareFn(m, v.v1, v.v2, v.matches)
			if m.fatal != v.fatal {
				t.Fatalf("%v != %v", m.fatal, v.fatal)
			}
		})
	}
}
