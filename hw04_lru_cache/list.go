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
	Value interface{} // значение
	Next  *listItem   // следующий элемент
	Prev  *listItem   // предыдущий элемент
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

	second := l.First

	l.First = elem

	if second != nil {
		second.Prev = l.First
		l.First.Next = second
	}

	if l.Len() == 0 {
		l.Last = l.First
	}

	l.Length++

	return l.First
}

func (l *list) PushBack(v interface{}) *listItem {
	elem := &listItem{Value: v, Next: nil, Prev: l.Last}
	previous := l.Last
	l.Last = elem

	if previous != nil {
		previous.Next = l.Last
	}

	if l.Len() == 0 {
		l.First = l.Last
	}

	l.Length++

	return l.Last
}

func (l *list) Remove(i *listItem) {
	previous := i.Prev
	next := i.Next

	i.Next = nil
	i.Prev = nil

	if previous != nil {
		previous.Next = next
	} else {
		l.First = next
	}

	if next != nil {
		next.Prev = previous
	} else {
		l.Last = previous
	}

	l.Length--
}

func (l *list) MoveToFront(i *listItem) {
	previous := i.Prev
	if previous == nil {
		// no need to move first element to front
		return
	}

	next := i.Next
	second := l.First

	previous.Next = next

	if next == nil {
		l.Last = previous
	} else {
		next.Prev = previous
	}

	i.Next = second
	second.Prev = i
	i.Prev = nil
	l.First = i
}

func NewList() List {
	return &list{}
}
