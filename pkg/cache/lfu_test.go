package cache_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/cache"
)

func newLFUCache[K comparable, V any](cap int, ttl time.Duration) cache.Cache[K, V] {
	return cache.New[K, V](cache.Config{
		Capacity:  cap,
		TTL:       ttl,
		TypeCache: cache.LFU,
	})
}

func TestLFU_SetGet(t *testing.T) {
	c := newLFUCache[string, int](3, 0)

	c.Set("1", 1)
	c.Set("2", 2)

	v, ok := c.Get("1")
	require.True(t, ok)
	require.Equal(t, 1, v)

	v, ok = c.Get("2")
	require.True(t, ok)
	require.Equal(t, 2, v)
}

func TestLFU_EvictionByFreq(t *testing.T) {
	c := newLFUCache[string, int](2, 0)

	c.Set("1", 10)
	c.Set("2", 20)
	c.Get("1")

	c.Set("3", 30)

	_, ok := c.Get("2")
	require.False(t, ok)

	v, ok := c.Get("1")
	require.True(t, ok)
	require.Equal(t, 10, v)

	v, ok = c.Get("3")
	require.True(t, ok)
	require.Equal(t, 30, v)
}

func TestLFU_TTLExpiration(t *testing.T) {
	c := newLFUCache[string, int](2, 50*time.Millisecond)

	c.Set("1", 1)
	time.Sleep(60 * time.Millisecond)

	_, ok := c.Get("1")
	require.False(t, ok)
}

func TestLFU_Delete(t *testing.T) {
	c := newLFUCache[string, int](2, 0)

	c.Set("123", 100)
	c.Delete("123")

	_, ok := c.Get("123")
	require.False(t, ok)
}

func TestLFU_Clear(t *testing.T) {
	c := newLFUCache[string, int](3, 0)

	c.Set("1", 1)
	c.Set("2", 2)
	c.Clear()

	_, ok := c.Get("1")
	require.False(t, ok)
	_, ok = c.Get("2")
	require.False(t, ok)
}
