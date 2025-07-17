package cache_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/cache"
)

func newLRUCache[K comparable, V any](capacity int, ttl time.Duration) cache.Cache[K, V] {
	return cache.New[K, V](cache.Config{
		Capacity:  capacity,
		TTL:       ttl,
		TypeCache: cache.LRU,
	})
}

func TestLRU_SetGet(t *testing.T) {
	c := newLRUCache[string, int](2, 0)

	c.Set("1", 1)
	c.Set("2", 2)

	v, ok := c.Get("1")
	require.True(t, ok)
	require.Equal(t, 1, v)

	v, ok = c.Get("2")
	require.True(t, ok)
	require.Equal(t, 2, v)
}

func TestLRU_Eviction(t *testing.T) {
	c := newLRUCache[string, int](2, 0)

	c.Set("1", 1)
	c.Set("2", 2)
	c.Get("1")
	c.Set("3", 3)

	_, ok := c.Get("2")
	require.False(t, ok)

	v, ok := c.Get("1")
	require.True(t, ok)
	require.Equal(t, 1, v)

	v, ok = c.Get("3")
	require.True(t, ok)
	require.Equal(t, 3, v)
}

func TestLRU_TTLExpiration(t *testing.T) {
	c := newLRUCache[string, int](2, 50*time.Millisecond)

	c.Set("x", 100)
	time.Sleep(100 * time.Millisecond)

	_, ok := c.Get("x")
	require.False(t, ok)
}

func TestLRU_DeleteClear(t *testing.T) {
	c := newLRUCache[string, int](2, 0)

	c.Set("1", 10)
	c.Delete("1")

	_, ok := c.Get("1")
	require.False(t, ok)

	c.Set("2", 20)
	c.Clear()

	_, ok = c.Get("2")
	require.False(t, ok)
}
