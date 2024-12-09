package main

import "sync"

// KV represents a thread-safe key-value store
type KV struct {
	mu   sync.RWMutex
	data map[string][]byte
}

// NewKV initializes a new key-value store
func NewKV() *KV {
	return &KV{
		data: make(map[string][]byte),
	}
}

// Set adds or updates a key-value pair
func (kv *KV) Set(key, val []byte) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.data[string(key)] = val
	return nil
}

// Get retrieves the value for a key
func (kv *KV) Get(key []byte) ([]byte, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	val, ok := kv.data[string(key)]
	return val, ok
}
