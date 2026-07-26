package compare

import (
	"reflect"
)

// TCompatible represents a type compatible with Helper and Fatalf from
// testing.T. This allows T itself to be mocked
type TCompatible interface {
	Helper()
	Fatalf(string, ...any)
}

// Compare compares two values of the same type for equality and fails the
// current test if they are not equal.
func Compare[T comparable](t TCompatible, v1, v2 T, matches bool) {
	t.Helper()
	if matches {
		if v1 != v2 {
			t.Fatalf("%v != %v", v1, v2)
		}
	} else {
		if v1 == v2 {
			t.Fatalf("%v == %v", v1, v2)
		}
	}
}

// CompareFn compares two functions for equality. Note that this doesn't take
// things like inline optimization or captured variables into account and
// therefore cannot compare for perfect uniqueness. However, it should work
// well enough for simple things like testing variables pointing to functions.
func CompareFn[T any](t TCompatible, v1, v2 T, matches bool) {
	t.Helper()
	if reflect.TypeOf(v1).Kind() != reflect.Func ||
		reflect.TypeOf(v2).Kind() != reflect.Func {
		t.Fatalf("both arguments must be functions")
		return
	}
	Compare(
		t,
		uintptr(reflect.ValueOf(v1).UnsafePointer()),
		uintptr(reflect.ValueOf(v2).UnsafePointer()),
		matches,
	)
}
