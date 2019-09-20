package dlist

type node struct {
	prev  *node
	next  *node
	value interface{}
}

// DoubleLinkedList - realization list
type DoubleLinkedList struct {
	head *node
	tail *node
}

// New - create new list
func New() *DoubleLinkedList {
	return &DoubleLinkedList{}
}

// Len - lenght of list, counter of list elements
func (d *DoubleLinkedList) Len() int {
	count := 0
	cur := d.head
	for ; cur != nil; cur = cur.next {
		count++
	}
	return count
}

// PopFront - get and pop first element
func (d *DoubleLinkedList) PopFront() interface{} {
	node := d.head
	return d.detach(node)
}

// PopBack - get and pop last element
func (d *DoubleLinkedList) PopBack() interface{} {
	node := d.tail
	return d.detach(node)
}

func (d *DoubleLinkedList) detach(n *node) interface{} {
	if n == nil {
		return nil
	}
	if n == d.head {
		d.head = n.next
	}
	if n == d.tail {
		d.tail = n.prev
	}
	if n.prev != nil {
		n.prev.next = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	}
	return n.value
}

func (d *DoubleLinkedList) new(v interface{}) *node {
	n := &node{value: v}
	if d.head == nil {
		d.head = n
		d.tail = n
	}
	return n
}

// PushFront - push element at begin of list
func (d *DoubleLinkedList) PushFront(v interface{}) {
	n := d.new(v)
	if d.head != n {
		d.head.prev = n
		n.next = d.head
		d.head = n
	}
}

// PushBack - push element at end of list
func (d *DoubleLinkedList) PushBack(v interface{}) {
	n := d.new(v)
	if d.tail != n {
		d.tail.next = n
		n.prev = d.tail
		d.tail = n
	}
}

// Remove - remove element from list by index
func (d *DoubleLinkedList) Remove(index int) {
	node := d.find(index)
	d.detach(node)
}

// Get - get element by index
func (d *DoubleLinkedList) Get(index int) interface{} {
	node := d.find(index)
	if node != nil {
		return node.value
	}
	return nil
}

func (d *DoubleLinkedList) find(index int) *node {
	if index < 0 {
		return nil
	}
	cur := d.head
	for ; index > 0 && cur != nil; index-- {
		cur = cur.next
	}
	return cur
}
