package cache

import (
	"sync"
	"time"
)

type lfuEnt[K comparable, V any] struct {
	key       K
	value     V
	freq      int
	expiresAt time.Time
}

type lfu[K comparable, V any] struct {
	mu       sync.Mutex
	capacity int
	ttl      time.Duration
	data     map[K]*lfuEnt[K, V]
	freqMap  map[int]map[K]*lfuEnt[K, V]
	minFreq  int
}

func newLFU[K comparable, V any](cfg Config) *lfu[K, V] {
	if cfg.Capacity <= 0 {
		panic("capacity must be positive")
	}
	return &lfu[K, V]{
		capacity: cfg.Capacity,
		ttl:      cfg.TTL,
		data:     make(map[K]*lfuEnt[K, V]),
		freqMap:  make(map[int]map[K]*lfuEnt[K, V]),
		minFreq:  0,
	}
}

func (c *lfu[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.data[key]
	if !ok {
		var zero V
		return zero, false
	}

	if c.ttl > 0 && time.Now().After(entry.expiresAt) {
		c.deleteEntry(entry)
		var zero V
		return zero, false
	}

	c.incrementFreq(entry)

	return entry.value, true
}

func (c *lfu[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, ok := c.data[key]; ok {
		entry.value = value
		entry.expiresAt = c.expiration()
		c.incrementFreq(entry)
		return
	}

	if len(c.data) >= c.capacity {
		c.evict()
	}

	entry := &lfuEnt[K, V]{
		key:       key,
		value:     value,
		freq:      1,
		expiresAt: c.expiration(),
	}
	c.data[key] = entry
	if c.freqMap[1] == nil {
		c.freqMap[1] = make(map[K]*lfuEnt[K, V])
	}
	c.freqMap[1][key] = entry
	c.minFreq = 1
}

func (c *lfu[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.data[key]
	if !ok {
		return
	}
	c.deleteEntry(entry)
}

func (c *lfu[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[K]*lfuEnt[K, V])
	c.freqMap = make(map[int]map[K]*lfuEnt[K, V])
	c.minFreq = 0
}

func (c *lfu[K, V]) incrementFreq(entry *lfuEnt[K, V]) {
	delete(c.freqMap[entry.freq], entry.key)
	if len(c.freqMap[entry.freq]) == 0 {
		delete(c.freqMap, entry.freq)
		if c.minFreq == entry.freq {
			c.minFreq++
		}
	}
	entry.freq++
	if c.freqMap[entry.freq] == nil {
		c.freqMap[entry.freq] = make(map[K]*lfuEnt[K, V])
	}
	c.freqMap[entry.freq][entry.key] = entry
}

func (c *lfu[K, V]) evict() {
	keys := c.freqMap[c.minFreq]
	for k, entry := range keys {
		c.deleteEntry(entry)
		delete(keys, k)
		break
	}
	if len(c.freqMap[c.minFreq]) == 0 {
		delete(c.freqMap, c.minFreq)
	}
}

func (c *lfu[K, V]) deleteEntry(entry *lfuEnt[K, V]) {
	delete(c.data, entry.key)
	if keys, ok := c.freqMap[entry.freq]; ok {
		delete(keys, entry.key)
		if len(keys) == 0 {
			delete(c.freqMap, entry.freq)
			if c.minFreq == entry.freq {
				c.minFreq = 0
			}
		}
	}
}

func (c *lfu[K, V]) expiration() time.Time {
	if c.ttl == 0 {
		return time.Time{}
	}
	return time.Now().Add(c.ttl)
}
