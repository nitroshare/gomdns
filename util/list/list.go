package list

// Element represents an element of a List.
type Element[T any] struct {
	Value T
	Prev  *Element[T]
	Next  *Element[T]
}

// List provides a bare-bones implementation of a generic doubly-linked list.
type List[T any] struct {
	Back  *Element[T]
	Front *Element[T]
	Len   int
}

// New creates a new List instance.
func New[T any]() *List[T] {
	return &List[T]{}
}

// Add adds the specified value to the front of the list and returns it.
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

// Remove removes the element from the list and returns the next element.
func (l *List[T]) Remove(e *Element[T]) *Element[T] {
	if e.Prev != nil {
		e.Prev.Next = e.Next
	}
	if e.Next != nil {
		e.Next.Prev = e.Prev
	}
	l.Len--
	return e.Next
}
