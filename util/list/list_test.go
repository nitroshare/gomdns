package list

import (
	"testing"

	"github.com/nitroshare/gomdns/util/test"
)

const (
	value1 = "test1"
	value2 = "test2"
	value3 = "test3"
)

func TestList(t *testing.T) {
	l := &List[string]{}
	test.Compare(t, l.Len, 0)
	test.Compare(t, l.Back, nil)
	test.Compare(t, l.Front, nil)
	e1 := l.Add(value1)
	test.Compare(t, l.Len, 1)
	test.Compare(t, l.Back, e1)
	test.Compare(t, l.Front, e1)
	test.Compare(t, e1.Value, value1)
	test.Compare(t, e1.Prev, nil)
	test.Compare(t, e1.Next, nil)
	e2 := l.Add(value2)
	test.Compare(t, l.Len, 2)
	test.Compare(t, l.Back, e2)
	test.Compare(t, l.Front, e1)
	test.Compare(t, e2.Value, value2)
	test.Compare(t, e1.Prev, nil)
	test.Compare(t, e1.Next, e2)
	test.Compare(t, e2.Prev, e1)
	test.Compare(t, e2.Next, nil)
	e3 := l.Add(value3)
	test.Compare(t, l.Len, 3)
	er := l.Remove(e2)
	test.Compare(t, l.Len, 2)
	test.Compare(t, l.Back, e3)
	test.Compare(t, l.Front, e1)
	test.Compare(t, er, e3)
	test.Compare(t, e1.Prev, nil)
	test.Compare(t, e1.Next, e3)
	test.Compare(t, e3.Prev, e1)
	test.Compare(t, e3.Next, nil)
}
