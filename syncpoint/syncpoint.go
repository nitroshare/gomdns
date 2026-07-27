package syncpoint

// SyncPoint provides a way to easily synchronize execution between goroutines
// during testing. It is designed to do nothing (NO-OP) until activated, at
// which point, it will then block calls to Wait() until Trigger() is invoked.
type SyncPoint struct {

	// Wait will pause the current goroutine until Trigger is called. This is
	// a NO-OP until Activate is invoked.
	Wait func()

	// Trigger will resume the goroutine that invoked Wait.
	Trigger func()

	// TriggerChan returns a channel that will pause until Wait is called or
	// returns nil if Activate hasn't been invoked.
	TriggerChan func() chan<- any

	chanTrigger chan any
}

// New creates a new SyncPoint instance.
func New() *SyncPoint {
	s := &SyncPoint{
		chanTrigger: make(chan any),
	}
	s.Deactivate()
	return s
}

// Activate enables the SyncPoint.
func (s *SyncPoint) Activate() {
	s.Wait = func() {
		<-s.chanTrigger
	}
	s.Trigger = func() {
		s.chanTrigger <- nil
	}
	s.TriggerChan = func() chan<- any {
		return s.chanTrigger
	}
}

// Deactivate disables the SyncPoint.
func (s *SyncPoint) Deactivate() {
	s.Wait = func() {}
	s.Trigger = func() {}
	s.TriggerChan = func() chan<- any { return nil }
}
