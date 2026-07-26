package list

// Element represents an element of a List.
type Element[T any] struct {
	Value T
	Prev  *Element[T]
	Next  *Element[T]
}

// List provides a bare-bones implementation of a generic doubly-linked list.
// This type is not safe for use in multiple goroutines.
type List[T any] struct {
	Back  *Element[T]
	Front *Element[T]
	Len   int
}

// Add adds the specified value to the back of the list and returns it.
func (l *List[T]) Add(v T) *Element[T] {
	e := &Element[T]{
		Value: v,
		Prev:  l.Back,
	}
	if l.Back != nil {
		l.Back.Next = e
	}
	if l.Front == nil {
		l.Front = e
	}
	l.Back = e
	l.Len++
	return e
}

// Remove removes the specified element from the list and returns the next
// one (if applicable).
func (l *List[T]) Remove(e *Element[T]) *Element[T] {
	if e.Prev != nil {
		e.Prev.Next = e.Next
	} else {
		l.Front = e.Next
	}
	if e.Next != nil {
		e.Next.Prev = e.Prev
	} else {
		l.Back = e.Prev
	}
	l.Len--
	return e.Next
}

// PopFront removes the first element and returns it. This will return nil if
// there are no elements in the list.
func (l *List[T]) PopFront() *Element[T] {
	e := l.Front
	if e != nil {
		l.Remove(e)
	}
	return e
}

// PopBack removes the last element and returns it. This will return nil if
// there are no elements in the list.
func (l *List[T]) PopBack() *Element[T] {
	e := l.Back
	if e != nil {
		l.Remove(e)
	}
	return e
}
