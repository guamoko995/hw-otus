package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  l.front,
	}
	if l.front != nil { // случай не пустого списка
		l.front.Prev = newItem
		l.front = newItem
		l.len++
	} else { // случай пустого списка
		l.front = newItem
		l.back = newItem
		l.len = 1
	}
	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Prev:  l.back,
	}
	if l.back != nil { // случай не пустого списка
		l.back.Next = newItem
		l.back = newItem
		l.len++
	} else { // случай пустого списка
		l.front = newItem
		l.back = newItem
		l.len = 1
	}
	return newItem
}

func (l *list) Remove(i *ListItem) {
	if l.len == 1 {
		l.len = 0
		l.front = nil
		l.back = nil
		return
	}
	switch i {
	case l.front:
		l.front = i.Next
		l.front.Prev = nil
	case l.back:
		l.back = i.Prev
		l.back.Next = nil
	default:
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.front {
		return
	}
	if i == l.back {
		l.back = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	i.Prev.Next = i.Next
	i.Prev = nil
	i.Next = l.front
	l.front.Prev = i
	l.front = i
}

func NewList() List {
	return new(list)
}
