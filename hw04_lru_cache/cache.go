package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()
	if item, ok := l.items[key]; ok {
		l.queue.MoveToFront(item)
		l.queue.Front().Value = cacheItem{
			key:   key,
			value: value,
		}
		return true
	}
	if l.queue.Len() == l.capacity {
		delete(l.items, l.queue.Back().Value.(cacheItem).key)
		l.queue.Remove(l.queue.Back())
	}
	l.items[key] = l.queue.PushFront(cacheItem{
		key:   key,
		value: value,
	})
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()
	if item, ok := l.items[key]; ok {
		l.queue.MoveToFront(item)
		return item.Value.(cacheItem).value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.Mutex.Lock()
	*l = lruCache{
		capacity: l.capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, l.capacity),
	}
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
