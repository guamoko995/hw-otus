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
	newList := &ListItem{
		Value: v,
		Next:  l.front,
	}
	if l.front != nil { // случай не пустого списка
		l.front.Prev = newList
		l.front = newList
		l.len++
	} else { // случай пустого списка
		l.front = newList
		l.back = newList
		l.len = 1
	}
	return newList
}

func (l *list) PushBack(v interface{}) *ListItem {
	newList := &ListItem{
		Value: v,
		Prev:  l.back,
	}
	if l.back != nil { // случай не пустого списка
		l.back.Next = newList
		l.back = newList
		l.len++
	} else { // случай пустого списка
		l.front = newList
		l.back = newList
		l.len = 1
	}
	return newList
}

func (l *list) Remove(i *ListItem) {
	switch i {
	case l.front:
		l.front = i.Next
		if l.front != nil {
			l.front.Prev = nil
			l.len--
		} else {
			l.back = nil
			l.len = 0
		}
	case l.back:
		l.back = i.Prev
		if l.back != nil {
			l.back.Next = nil
			l.len--
		} else {
			l.front = nil
			l.len = 0
		}
	default:
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
		l.len--
	}
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
