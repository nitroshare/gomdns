package list

import (
	"testing"

	"github.com/nitroshare/gomdns/util/compare"
)

const (
	value1 = "test1"
	value2 = "test2"
	value3 = "test3"
)

func TestList(t *testing.T) {
	l := &List[string]{}
	compare.Compare(t, l.Len, 0, true)
	compare.Compare(t, l.Back, nil, true)
	compare.Compare(t, l.Front, nil, true)
	e1 := l.Add(value1)
	compare.Compare(t, l.Len, 1, true)
	compare.Compare(t, l.Back, e1, true)
	compare.Compare(t, l.Front, e1, true)
	compare.Compare(t, e1.Value, value1, true)
	compare.Compare(t, e1.Prev, nil, true)
	compare.Compare(t, e1.Next, nil, true)
	e2 := l.Add(value2)
	compare.Compare(t, l.Len, 2, true)
	compare.Compare(t, l.Back, e2, true)
	compare.Compare(t, l.Front, e1, true)
	compare.Compare(t, e2.Value, value2, true)
	compare.Compare(t, e1.Prev, nil, true)
	compare.Compare(t, e1.Next, e2, true)
	compare.Compare(t, e2.Prev, e1, true)
	compare.Compare(t, e2.Next, nil, true)
	e3 := l.Add(value3)
	compare.Compare(t, l.Len, 3, true)
	er := l.Remove(e2)
	compare.Compare(t, l.Len, 2, true)
	compare.Compare(t, l.Back, e3, true)
	compare.Compare(t, l.Front, e1, true)
	compare.Compare(t, er, e3, true)
	compare.Compare(t, e1.Prev, nil, true)
	compare.Compare(t, e1.Next, e3, true)
	compare.Compare(t, e3.Prev, e1, true)
	compare.Compare(t, e3.Next, nil, true)
	l.Remove(e1)
	l.Remove(e3)
	compare.Compare(t, l.Len, 0, true)
	compare.Compare(t, l.Back, nil, true)
	compare.Compare(t, l.Front, nil, true)
}

func TestPop(t *testing.T) {
	var (
		l  = &List[string]{}
		e1 = l.Add(value1)
		e2 = l.Add(value1)
	)
	ef := l.PopFront()
	compare.Compare(t, ef, e1, true)
	eb := l.PopBack()
	compare.Compare(t, eb, e2, true)
}
