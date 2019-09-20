package dlist

import "testing"

func check(t *testing.T, list *DoubleLinkedList, len int, values []int) {
	if list.Len() != len {
		t.Helper()
		t.Errorf("[%s] not valid list len. Got %d, was expected %d", t.Name(), list.Len(), len)
		return
	}
	for i := 0; i < len; i++ {
		if !checkValue(t, list.Get(i), values[i]) {
			break
		}
	}
}

func checkValue(t *testing.T, v interface{}, value int) bool {
	intvalue, ok := v.(int)
	if !ok {
		t.Helper()
		t.Errorf("[%s] returned value has not valid type", t.Name())
		return false
	}
	if intvalue != value {
		t.Helper()
		t.Errorf("[%s] returned value not valid. Got %d, was expected %d", t.Name(), intvalue, value)
		return false
	}
	return true
}

func checkNilValue(t *testing.T, v interface{}) {
	if v != nil {
		t.Helper()
		t.Errorf("[%s] returned not Nil value", t.Name())
	}
}

func TestZeroElements(t *testing.T) {
	l := New()
	check(t, l, 0, []int{})
	l.PopBack()
	check(t, l, 0, []int{})
	l.PopFront()
	check(t, l, 0, []int{})
	checkNilValue(t, l.Get(0))
	checkNilValue(t, l.Get(100))
}

func TestSingleElement(t *testing.T) {
	l := New()
	l.PushFront(1)
	check(t, l, 1, []int{1})
	checkValue(t, l.PopBack(), 1)
	check(t, l, 0, []int{})
	l.PushBack(2)
	check(t, l, 1, []int{2})
	checkValue(t, l.PopFront(), 2)
	check(t, l, 0, []int{})
}

func TestMuliplyElements(t *testing.T) {
	l := New()
	for i := 0; i < 10; i++ {
		l.PushBack(i)
	}
	check(t, l, 10, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	l2 := New()
	for i := 0; i < 10; i++ {
		l2.PushFront(i)
	}
	check(t, l2, 10, []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0})

	l.Remove(5)
	check(t, l, 9, []int{0, 1, 2, 3, 4, 6, 7, 8, 9})
	l.PopFront()
	l.PopBack()
	check(t, l, 7, []int{1, 2, 3, 4, 6, 7, 8})
	checkNilValue(t, l.Get(100))
}
