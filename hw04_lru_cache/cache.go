package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

type cacheItem struct {
	key   Key
	value any
}

func (c *lruCache) Set(key Key, value any) bool {
	if item, found := c.items[key]; found {
		item.Value.(*cacheItem).value = value
		c.queue.MoveToFront(item)
		return true
	}

	item := c.queue.PushFront(&cacheItem{key, value})
	c.items[key] = item

	if c.queue.Len() > c.capacity {
		last := c.queue.Back()
		if last != nil {
			c.queue.Remove(last)
			delete(c.items, last.Value.(*cacheItem).key)
		}
	}
	return false
}

func (c *lruCache) Get(key Key) (any, bool) {
	if item, found := c.items[key]; found {
		c.queue.MoveToFront(item)
		return item.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem)
}
