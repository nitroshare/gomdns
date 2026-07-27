package syncpoint

import (
	"testing"
)

func TestSyncPoint(t *testing.T) {
	s := New()
	s.Wait()
	s.Trigger()
	s.TriggerChan()
	s.Activate()
	select {
	case s.TriggerChan() <- nil:
		t.Fatal("shouldn't be able to send on channel")
	default:
	}
	go func() {
		s.Trigger()
	}()
	s.Wait()
}
