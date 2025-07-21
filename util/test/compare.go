package test

import (
	"testing"
)

// Compare compares two values of the same type for equality and fails the
// current test if they are not equal.
func Compare[T comparable](t *testing.T, v1, v2 T) {
	if v1 != v2 {
		t.Fatalf("\"%v\" != \"%v\"", v1, v2)
	}
}
