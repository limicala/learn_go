package cache

import (
	"container/list"
)

type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

type Cache struct {
	size, capacity int64      // capacity=0 meaning no capacity limit
	list           *list.List // double linked list
	cache          map[string]*list.Element

	OnRemove func(key string, value Value)
}

func getEntrySize(entry *entry) int {
	return len(entry.key) + entry.value.Len()
}

func New(capacity int64, onRemoveCallback func(string, Value)) *Cache {
	return &Cache{
		size:     0,
		capacity: capacity,
		list:     list.New(),
		cache:    make(map[string]*list.Element),

		OnRemove: onRemoveCallback,
	}
}

func (c *Cache) Len() int {
	return c.list.Len()
}

func (c *Cache) Size() int64 {
	return c.capacity
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := c.cache[key]; ok {
		c.list.MoveToFront(element)
		entry := element.Value.(*entry)
		value = entry.value
	}
	return
}

func (c *Cache) Add(key string, value Value) {
	if element, ok := c.cache[key]; ok {
		c.list.MoveToFront(element)
		entry := element.Value.(*entry)
		oldSize := getEntrySize(entry)
		entry.value = value
		c.size += int64(getEntrySize(entry) - oldSize)
	} else {
		entry := &entry{key, value}
		element := c.list.PushFront(entry)
		c.cache[key] = element
		c.size += int64(getEntrySize(entry))
	}
	for c.capacity != 0 && c.size > c.capacity {
		c.RemoveOldest()
	}
}

func (c *Cache) RemoveOldest() {
	if element := c.list.Back(); element != nil {
		c.list.Remove(element)
		// Type assertions https://go.dev/tour/methods/15
		entry := element.Value.(*entry)
		delete(c.cache, entry.key)
		c.size -= int64(getEntrySize(entry))
		if c.OnRemove != nil {
			c.OnRemove(entry.key, entry.value)
		}
	}
}
