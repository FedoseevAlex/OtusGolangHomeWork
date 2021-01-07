package hw04_lru_cache //nolint:golint,stylecheck

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	Cap   int
	Queue List
	Items map[Key]*listItem
}

type cacheItem struct {
	Key   Key
	Value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) (ok bool) {
	item, ok := c.Items[key]

	if ok {
		item.Value = cacheItem{Key: key, Value: value}
		c.Queue.MoveToFront(item)
	} else {
		c.Queue.PushFront(cacheItem{Key: key, Value: value})
		c.Items[key] = c.Queue.Front()
	}

	if c.Queue.Len() > c.Cap {
		lastItem := c.Queue.Back()
		c.Queue.Remove(lastItem)
		delete(c.Items, lastItem.Value.(cacheItem).Key)
	}
	return
}

func (c *lruCache) Get(key Key) (value interface{}, ok bool) {
	item, ok := c.Items[key]
	if ok {
		c.Queue.MoveToFront(item)
		value = item.Value.(cacheItem).Value
	}
	return
}

func (c *lruCache) Clear() {
	c.Queue = NewList()
	c.Items = make(map[Key]*listItem)
}

func NewCache(capacity int) Cache {
	return &lruCache{Cap: capacity, Queue: NewList(), Items: make(map[Key]*listItem)}
}
