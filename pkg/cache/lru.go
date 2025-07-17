package cache

import (
	"container/list"
	"sync"
	"time"
)

type lruEnt[K comparable, V any] struct {
	key       K
	value     V
	expiresAt time.Time
}

type lru[K comparable, V any] struct {
	mu       sync.Mutex
	capacity int
	ttl      time.Duration
	data     map[K]*list.Element
	evict    *list.List
}

// newLRU создает новый экземпляр LRU алгоритма кэширования.
func newLRU[K comparable, V any](cfg Config) *lru[K, V] {
	if cfg.Capacity <= 0 {
		panic("capacity must be positive")
	}
	return &lru[K, V]{
		capacity: cfg.Capacity,
		ttl:      cfg.TTL,
		data:     make(map[K]*list.Element),
		evict:    list.New(),
	}
}

// Get возвращает значение по ключу.
func (c *lru[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.data[key]
	if !ok {
		var zero V
		return zero, false
	}

	e := elem.Value.(lruEnt[K, V])
	if c.ttl > 0 && time.Now().After(e.expiresAt) {
		c.removeElement(elem)
		var zero V
		return zero, false
	}

	c.evict.MoveToBack(elem)
	return e.value, true
}

// Set добавляет значение в кэш.
func (c *lru[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.data[key]; ok {
		elem.Value = lruEnt[K, V]{key, value, c.exp()}
		c.evict.MoveToBack(elem)
		return
	}

	if c.evict.Len() >= c.capacity {
		c.evictOldest()
	}

	elem := c.evict.PushBack(lruEnt[K, V]{key, value, c.exp()})
	c.data[key] = elem
}

// Delete удаляет ключ, если он существует.
func (c *lru[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.data[key]; ok {
		c.removeElement(elem)
	}
}

// Clear удаляет все элементы из кэша.
func (c *lru[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[K]*list.Element)
	c.evict.Init()
}

// Clear очистка кэша
func (c *lru[K, V]) evictOldest() {
	oldest := c.evict.Front()
	if oldest != nil {
		c.removeElement(oldest)
	}
}

func (c *lru[K, V]) removeElement(elem *list.Element) {
	e := elem.Value.(lruEnt[K, V])
	delete(c.data, e.key)
	c.evict.Remove(elem)
}

func (c *lru[K, V]) exp() time.Time {
	if c.ttl == 0 {
		return time.Time{}
	}
	return time.Now().Add(c.ttl)
}
