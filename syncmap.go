package goutil

import "sync"

type SyncMap[K Hashable, V any] struct {
	value map[K]V
	hasVal map[K]bool
	mu sync.Mutex
	null V
}

// NewMap creates a new synchronized map that uses sync.Mutex behind the scenes
func NewMap[K Hashable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		value: map[K]V{},
		hasVal: map[K]bool{},
	}
}

// Get returns a value or an error if it exists
func (syncmap *SyncMap[K, V]) Get(key K) (V, bool) {
	syncmap.mu.Lock()
	defer syncmap.mu.Unlock()

	if hasVal, ok := syncmap.hasVal[key]; !ok || !hasVal {
		return syncmap.null, false
	}else if val, ok := syncmap.value[key]; ok {
		return val, true
	}

	return syncmap.null, false
}

// Set sets or adds a new key with a value
func (syncmap *SyncMap[K, V]) Set(key K, value V) {
	syncmap.mu.Lock()
	defer syncmap.mu.Unlock()

	syncmap.value[key] = value
	syncmap.hasVal[key] = true
}

// Del removes an item by key
func (syncmap *SyncMap[K, V]) Del(key K){
	syncmap.mu.Lock()
	defer syncmap.mu.Unlock()

	delete(syncmap.value, key)
	delete(syncmap.hasVal, key)
}

// Has returns true if a key value exists in the list
func (syncmap *SyncMap[K, V]) Has(key K) bool {
	syncmap.mu.Lock()
	defer syncmap.mu.Unlock()

	if hasVal, ok := syncmap.hasVal[key]; !ok || !hasVal {
		return false
	}else if _, ok := syncmap.value[key]; ok {
		return true
	}

	return false
}

// ForEach runs a callback function for each key value pair
//
// in the callback, return true to continue, and false to break the loop
func (syncmap *SyncMap[K, V]) ForEach(cb func(key K, value V) bool){
	syncmap.mu.Lock()
	keyList := []K{}
	for key := range syncmap.value {
		keyList = append(keyList, key)
	}
	syncmap.mu.Unlock()
	
	for _, key := range keyList {
		syncmap.mu.Lock()

		if hasVal, ok := syncmap.hasVal[key]; !ok || !hasVal {
			delete(syncmap.value, key)
			delete(syncmap.hasVal, key)
			syncmap.mu.Unlock()
			continue
		}

		var val V
		if v, ok := syncmap.value[key]; ok {
			val = v
		}

		syncmap.mu.Unlock()

		if !cb(key, val) {
			break
		}
	}
}
