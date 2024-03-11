package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	mutex    sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, ok := c.items[key]
	newItem := &ListItem{Value: cacheItem{key: key, value: value}, Next: c.queue.Front(), Prev: nil}
	c.items[key] = newItem

	if !ok {
		if c.queue.Len() == c.capacity {
			c.queue.Remove(c.queue.Back())
			if backItemValue, ok := c.queue.Back().Value.(cacheItem); ok {
				delete(c.items, backItemValue.key)
			}
		}

		c.queue.PushFront(newItem)
	} else {
		c.queue.MoveToFront(newItem)
	}

	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, ok := c.items[key]

	if !ok {
		return nil, false
	}

	c.queue.MoveToFront(item)
	return item.Value.(cacheItem).value, ok
}

func (c *lruCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.queue.Front().Next = nil
	c.queue.Back().Prev = nil
	c.items = make(map[Key]*ListItem)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		mutex:    sync.Mutex{},
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
