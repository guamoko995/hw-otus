package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// newListWhithElements cоздает лист с указанными элементами.
func newListWhithElements(elements ...interface{}) List {
	l := NewList()
	for _, el := range elements {
		l.PushBack(el)
	}
	return l
}

// testListOfOneElement проверяет лист l на соответствие
// двусвязному списку из указанных элементов.
// Проверяются все ссылки и все значения.
func testListWhithElements(t *testing.T, l List, elements ...interface{}) {
	require.Equal(t, len(elements), l.Len())
	list := l.Front()
	for i := range elements {
		require.Equal(t, elements[i], list.Value)
		if i == 0 {
			require.Nil(t, list.Prev)
		} else {
			require.Equal(t, list.Prev.Next, list)
		}
		if i == len(elements)-1 {
			require.Nil(t, list.Next)
		} else {
			require.Equal(t, list.Next.Prev, list)
		}
		list = list.Next
	}
}

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})
	t.Run("PushFront", func(t *testing.T) {
		l := NewList()
		l.PushFront(1)
		testListWhithElements(t, l, 1)
		l.PushFront(2.0)
		testListWhithElements(t, l, 2.0, 1)
		l.PushFront("3")
		testListWhithElements(t, l, "3", 2.0, 1)
	})
	t.Run("PushBack", func(t *testing.T) {
		l := NewList()
		l.PushBack(1)
		testListWhithElements(t, l, 1)
		l.PushBack(2.0)
		testListWhithElements(t, l, 1, 2.0)
		l.PushBack("3")
		testListWhithElements(t, l, 1, 2.0, "3")
	})
	t.Run("Remove one", func(t *testing.T) {
		l := newListWhithElements("Один")
		l.Remove(l.Front())
		testListWhithElements(t, l)
	})
	t.Run("Remove front", func(t *testing.T) {
		l := newListWhithElements(1, 2, 3)
		l.Remove(l.Front())
		testListWhithElements(t, l, 2, 3)
	})
	t.Run("Remove back", func(t *testing.T) {
		l := newListWhithElements(1, 2, 3)
		l.Remove(l.Back())
		testListWhithElements(t, l, 1, 2)
	})
	t.Run("Remove", func(t *testing.T) {
		l := newListWhithElements(1, 2, 3)
		l.Remove(l.Front().Next)
		testListWhithElements(t, l, 1, 3)
	})
	t.Run("MoveToFront one", func(t *testing.T) {
		l := newListWhithElements(1)
		l.MoveToFront(l.Front())
		testListWhithElements(t, l, 1)
	})
	t.Run("MoveToFront front", func(t *testing.T) {
		l := newListWhithElements(1, 2, 3)
		l.MoveToFront(l.Front())
		testListWhithElements(t, l, 1, 2, 3)
	})
	t.Run("MoveToFront back", func(t *testing.T) {
		l := newListWhithElements(1, 2, 3)
		l.MoveToFront(l.Back())
		testListWhithElements(t, l, 3, 1, 2)
	})
	t.Run("MoveToFront", func(t *testing.T) {
		l := newListWhithElements(1, 2, 3)
		l.MoveToFront(l.Back().Prev)
		testListWhithElements(t, l, 2, 1, 3)
	})
	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}
