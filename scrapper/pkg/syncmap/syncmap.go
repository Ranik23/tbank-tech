package syncmap

import (
	"sync"
)

type SyncMap[K comparable, V any] struct {
	mu   sync.Mutex
	data map[K]V
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		data: make(map[K]V),
	}
}

func (m *SyncMap[K, V]) Store(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

func (m *SyncMap[K, V]) Load(key K) (value V, ok bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	value, ok = m.data[key]
	return
}

func (m *SyncMap[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

func (m *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}