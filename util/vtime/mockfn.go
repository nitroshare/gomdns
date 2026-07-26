package vtime

import (
	"time"
)

// MockNow returns the current mocked time. Although this can be set by
// reassigning Now, this is typically handled automatically by Mock.
func MockNow() time.Time {
	loop.chanNow <- nil
	return <-loop.chanTime
}

// MockAfter returns a channel that will send when the specified duration
// elapses via calls to Set or Advance.
func MockAfter(d time.Duration) <-chan time.Time {
	if loop.chanTest != nil {
		<-loop.chanTest
	}
	loop.chanAfter <- d
	return <-loop.chanTimeChan
}
