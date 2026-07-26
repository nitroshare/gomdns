package broadcaster

import (
	"sync"
)

// Broadcaster provides an implementation of the fan-out pattern, providing a
// single channel for sending and creating a new channel for each listener
// that subscribes.
type Broadcaster[T any] struct {
	chanSubscribe      chan any
	chanSubscribeRet   chan (<-chan T)
	chanUnsubscribe    chan (<-chan T)
	chanUnsubscribeRet chan any
	chanSend           chan T
	chanClose          chan any
	chanClosed         chan any
}

// In order to avoid blocking the main loop when sending to each subscriber,
// the goroutine servicing each subscriber immediately stores each incoming
// value in its own queue and tries sending in the same select as the cases
// receiving new values and watching for shutdown; because the send is non-
// blocking, the main loop can send to each goroutine without blocking

func run[T any](chanRec <-chan T, chanSend chan<- T) {
	defer close(chanSend)
	queue := []T{}
	for {
		var (
			chanSendOrNil chan<- T
			nextV         T
		)
		if len(queue) > 0 {
			chanSendOrNil = chanSend
			nextV = queue[0]
		}
		select {
		case v, ok := <-chanRec:
			if !ok {
				return
			}
			queue = append(queue, v)
		case chanSendOrNil <- nextV:
			queue = queue[1:]
		}
	}
}

func (b *Broadcaster[T]) loop() {
	defer close(b.chanClosed)
	var (
		wg          sync.WaitGroup
		subscribers = map[<-chan T]chan T{}
	)
	defer wg.Wait()
	for {
		select {
		case <-b.chanSubscribe:
			var (
				chanSub = make(chan T)
				chanRun = make(chan T)
			)
			subscribers[chanSub] = chanRun
			wg.Add(1)
			go func() {
				defer wg.Done()
				run(chanRun, chanSub)
			}()
			b.chanSubscribeRet <- chanSub
		case chanSub := <-b.chanUnsubscribe:
			if chanRun, ok := subscribers[chanSub]; ok {
				close(chanRun)
				delete(subscribers, chanSub)
			}
			b.chanUnsubscribeRet <- nil
		case v, ok := <-b.chanSend:
			if !ok {
				for _, chanRun := range subscribers {
					close(chanRun)
				}
				return
			}
			for _, chanRun := range subscribers {
				chanRun <- v
			}
		}
	}
}

// New creates a new Broadcaster instance.
func New[T any]() *Broadcaster[T] {
	b := &Broadcaster[T]{
		chanSubscribe:      make(chan any),
		chanSubscribeRet:   make(chan (<-chan T)),
		chanUnsubscribe:    make(chan (<-chan T)),
		chanUnsubscribeRet: make(chan any),
		chanSend:           make(chan T),
		chanClosed:         make(chan any),
	}
	go b.loop()
	return b
}

// Subscribe returns a channel that will receive all broadcasts until either
// the broadcaster shuts down or Unsubscribe is called. Either will result in
// the returned channel being closed.
func (b *Broadcaster[T]) Subscribe() <-chan T {
	b.chanSubscribe <- nil
	return <-b.chanSubscribeRet
}

// Unsubscribe closes the channel returned from Subscribe.
func (b *Broadcaster[T]) Unsubscribe(chanSub <-chan T) {
	b.chanUnsubscribe <- chanSub
	<-b.chanUnsubscribeRet
}

// Send broadcasts the provided value to all subscribers. Note that this does
// not make a deep copy of the value.
func (b *Broadcaster[T]) Send(v T) {
	b.chanSend <- v
}

// Close shuts down the broadcaster. All subscribers currently listening for
// broadcasts will have their receiving channels closed.
func (b *Broadcaster[T]) Close() {
	close(b.chanSend)
	<-b.chanClosed
}
