// Package cache предоставляет простой in-memory кэш с поддержкой стратегий вытеснения LRU и LFU.
//
// Основные возможности:
//   - Ограничение по количеству элементов (Capacity)
//   - Автоматическая инвалидация по TTL
//   - Выбор стратегии вытеснения: LRU или LFU
//
// Пример:
//
//	cfg := cache.Config{
//	    Capacity:  100,
//	    TTL:       time.Minute,
//	    TypeCache: cache.LRU,
//	}
//	cache := cache.New[string, string](cfg)
//
// Инвалидация:
//   - Все элементы автоматически удаляются по истечении TTL.
//   - Устаревшие элементы очищаются при чтении (Get).
//
// Стратегии вытеснения:
//   - LRU (Least Recently Used)
//   - LFU (Least Frequently Used)
package cache

import (
	"time"
)

type Cache[K comparable, V any] interface {
	Get(key K) (V, bool)
	Set(key K, value V)
	Delete(key K)
	Clear()
}

type TypeCache int

const (
	LRU TypeCache = iota
	LFU
)

type Config struct {
	Capacity  int
	TTL       time.Duration
	TypeCache TypeCache
}

func New[K comparable, V any](cfg Config) Cache[K, V] {
	switch cfg.TypeCache {
	case LRU:
		return newLRU[K, V](cfg)
	case LFU:
		return newLFU[K, V](cfg)
	default:
		return nil
	}
}
