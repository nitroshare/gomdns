package vtime

import (
	"time"
)

var (
	// Now normally points to time.Now but can also be pointed to with
	// MockNow.
	Now func() time.Time

	// After normally points to time.After but can also be pointed to with
	// MockAfter.
	After func(d time.Duration) <-chan time.Time
)

func reset() {
	Now = time.Now
	After = time.After
}

func init() {
	reset()
}

// Mock replaces the time functions in this package with their mocked
// equivalents. Note that this is intended to be called at the beginning of a
// test and used with an associated defer Unmock() call. It is not safe for
// Mock or Unmock to be called from other goroutines, although the other
// mocked functions and types may be used anywhere while mocking is active.
func Mock() {
	loop = newMockLoop()
	Now = MockNow
	After = MockAfter
}

// Unmock replaces the mocked time functions with their original equivalents.
func Unmock() {
	loop.close()
	loop = nil
	reset()
}
