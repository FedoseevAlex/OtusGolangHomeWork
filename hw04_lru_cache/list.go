package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int
	Front() *listItem
	Back() *listItem
	PushFront(v interface{}) *listItem
	PushBack(v interface{}) *listItem
	Remove(i *listItem)
	MoveToFront(i *listItem)
}

type listItem struct {
	Value interface{}
	Next  *listItem
	Prev  *listItem
}

type list struct {
	First  *listItem
	Last   *listItem
	Length int
}

func (l *list) Len() int {
	return l.Length
}

func (l *list) Front() *listItem {
	return l.First
}

func (l *list) Back() *listItem {
	return l.Last
}

func (l *list) PushFront(v interface{}) *listItem {
	elem := &listItem{Value: v, Next: l.First, Prev: nil}

	// Second element after push to front
	second := l.First

	l.First = elem

	// Connect first element to second
	if second == nil {
		// If we pushed first element to front
		// then mark it as the last element too
		l.Last = l.First
	} else {
		second.Prev = l.First
		l.First.Next = second
	}

	l.Length++

	return l.First
}

func (l *list) PushBack(v interface{}) *listItem {
	elem := &listItem{Value: v, Next: nil, Prev: l.Last}

	previous := l.Last

	l.Last = elem

	// Connect last element to the end of list
	if previous == nil {
		// If we are here it means that we are
		// pushing first element to the back
		l.First = l.Last
	} else {
		previous.Next = l.Last
		l.Last.Prev = previous
	}

	l.Length++

	return l.Last
}

func (l *list) Remove(i *listItem) {
	// Save surrounding elements
	previous := i.Prev
	next := i.Next

	// Clear links in removable item
	i.Next = nil
	i.Prev = nil

	// Connect surrounding elements
	// to each other
	if previous == nil {
		// Removal of the first element
		l.First = next
	} else {
		previous.Next = next
	}

	if next == nil {
		// Removal of last element
		l.Last = previous
	} else {
		next.Prev = previous
	}

	l.Length--
}

func (l *list) MoveToFront(i *listItem) {
	previous := i.Prev
	if previous == nil {
		// If we are trying to move
		// first element to front then just return
		return
	}

	next := i.Next
	second := l.First

	// Connect next and previous elements
	previous.Next = next

	if next == nil {
		l.Last = previous
	} else {
		next.Prev = previous
	}

	// Connect specified and second elements
	i.Next = second
	second.Prev = i

	// Move specified element to front
	l.First = i
	i.Prev = nil
}

func NewList() List {
	return &list{}
}
