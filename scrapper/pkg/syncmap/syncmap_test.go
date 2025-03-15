package syncmap

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncMap_StoreAndLoad(t *testing.T) {
	// Создаем новый SyncMap
	sm := NewSyncMap[string, int]()

	// Добавляем элемент
	sm.Store("key1", 42)

	// Загружаем элемент
	value, ok := sm.Load("key1")
	require.True(t, ok, "Expected key to exist")
	assert.Equal(t, 42, value, "Expected value to be 42")
}

func TestSyncMap_Load_NotFound(t *testing.T) {
	// Создаем новый SyncMap
	sm := NewSyncMap[string, int]()

	// Загружаем элемент, которого нет в карте
	value, ok := sm.Load("nonexistent_key")
	assert.False(t, ok, "Expected key to not exist")
	assert.Equal(t, 0, value, "Expected default value for non-existent key")
}

func TestSyncMap_Delete(t *testing.T) {
	// Создаем новый SyncMap
	sm := NewSyncMap[string, int]()

	// Добавляем элемент
	sm.Store("key1", 42)

	// Удаляем элемент
	sm.Delete("key1")

	// Пытаемся загрузить удаленный элемент
	value, ok := sm.Load("key1")
	assert.False(t, ok, "Expected key to be deleted")
	assert.Equal(t, 0, value, "Expected default value after deletion")
}

func TestSyncMap_Range(t *testing.T) {
	// Создаем новый SyncMap
	sm := NewSyncMap[string, int]()

	// Добавляем элементы
	sm.Store("key1", 42)
	sm.Store("key2", 100)

	// Используем Range для обхода элементов
	var keys []string
	var values []int

	sm.Range(func(key string, value int) bool {
		keys = append(keys, key)
		values = append(values, value)
		// Прерываем после первого элемента
		return false
	})

	// Проверяем, что перебрали элементы
	assert.Len(t, keys, 1, "Expected to have one key")
	assert.Len(t, values, 1, "Expected to have one value")
	assert.Equal(t, "key1", keys[0], "Expected the first key to be 'key1'")
	assert.Equal(t, 42, values[0], "Expected the first value to be 42")
}

func TestSyncMap_Concurrent(t *testing.T) {
	sm := NewSyncMap[int, string]()

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sm.Store(i, "value")
		}()
	}


	wg.Wait()
	
	for i := 0; i < 1000; i++ {
		value, ok := sm.Load(i)
		assert.True(t, ok, "Expected key to exist")
		assert.Equal(t, "value", value, "Expected value to be 'value'")
	}
}

