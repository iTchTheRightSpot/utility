package cache

import (
	"github.com/iTchTheRightSpot/utility/log"
	"sync"
	"time"
)

type ICache[K any, V any] interface {
	Put(key K, value V)
	Get(key K) *V
	Delete(key K)
	Clear()
}

type InMemoryCache[K any, V any] struct {
	logger   log.ILogger
	cache    *sync.Map
	duration time.Duration
	size     int
}

type customValue[V any] struct {
	timer      *time.Timer
	value      V
	LastAccess time.Time
}

// SyncMapInMemoryCache duration is in minutes
func SyncMapInMemoryCache[K any, V any](l log.ILogger, duration, size int) ICache[K, V] {
	return &InMemoryCache[K, V]{
		logger:   l,
		cache:    &sync.Map{},
		duration: time.Duration(duration) * time.Minute * time.Second,
		size:     size,
	}
}

func (dep *InMemoryCache[K, V]) Length() int {
	var count int
	dep.cache.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

func (dep *InMemoryCache[K, V]) LeastUsed() K {
	var k K
	now := dep.logger.Date()
	lastAccess := &now
	dep.cache.Range(func(key, value any) bool {
		v := value.(customValue[V])
		if v.LastAccess.Before(*lastAccess) {
			k = key.(K)
			lastAccess = &v.LastAccess
		}
		return true
	})
	return k
}

func (dep *InMemoryCache[K, V]) Put(key K, value V) {
	length := dep.Length()
	if length == dep.size {
		k := dep.LeastUsed()
		dep.cache.Delete(k)
	}

	if _, ok := dep.cache.Load(key); ok {
		dep.cache.Delete(key)
	}

	timeout := time.AfterFunc(dep.duration, func() { dep.cache.Delete(key) })
	dep.cache.Store(key, customValue[V]{timer: timeout, value: value})
}

func (dep *InMemoryCache[K, V]) Get(key K) *V {
	value, ok := dep.cache.Load(key)
	if !ok {
		return nil
	}
	v := value.(customValue[V])
	v.LastAccess = dep.logger.Date()
	return &v.value
}

func (dep *InMemoryCache[K, V]) Delete(key K) {
	dep.cache.Delete(key)
}

func (dep *InMemoryCache[K, V]) Clear() {
	dep.cache.Range(func(key, value any) bool {
		dep.cache.Delete(key)
		return true
	})
}
