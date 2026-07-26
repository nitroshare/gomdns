package vtime

import (
	"testing"
	"time"

	"github.com/nitroshare/compare"
)

func TestMockUnmock(t *testing.T) {
	compare.CompareFn(t, Now, time.Now, true)
	compare.CompareFn(t, After, time.After, true)
	Mock()
	compare.CompareFn(t, Now, MockNow, true)
	compare.CompareFn(t, After, MockAfter, true)
	Unmock()
	compare.CompareFn(t, Now, time.Now, true)
	compare.CompareFn(t, After, time.After, true)
}
