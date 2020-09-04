package qqcache

import (
	"errors"
	"sync"
	"time"
)

const defaultEvictionInterval = 60 * time.Second

var (
	ErrWrongTypeIndex = errors.New("wrong type of the value to get by index")
	ErrWrongTypeLPush = errors.New("wrong type of the value to push list value")
	ErrWrongTypeHSet  = errors.New("wrong type of the value to set hash map value")
	ErrWrongTypeHGet  = errors.New("wrong type of the value to get hash map key")
	ErrNotFound       = errors.New("not value found by key")
)

// Opts represents the options to create new instance of Cache.
type Opts struct {
	// EvictionInterval is how often cache-cleaner will delete expired keys.
	EvictionInterval time.Duration
}

// Cache represents cache container.
type Cache struct {
	mux              sync.RWMutex
	data             map[string]entity
	evictionInterval time.Duration
	stopCleaner      chan struct{}
}

// New returns new instance of Cache.
// If eviction interval is equal or less that 0 - default eviction will be used.
func New(opts Opts) *Cache {
	c := &Cache{
		mux:              sync.RWMutex{},
		data:             make(map[string]entity),
		evictionInterval: opts.EvictionInterval,
		stopCleaner:      make(chan struct{}),
	}

	// Run cache cleaner
	go c.cacheCleaner()

	return c
}

// Shutdown stops cache cleaner.
func (c *Cache) Shutdown() {
	c.stopCleaner <- struct{}{}
}

// Set method sets value to cache by key with specific TTL.
// If given TTL <=0 then the key will never be expired.
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.data[key] = entity{
		value:        value,
		expiredAfter: validateExpiredAfter(ttl),
	}
}

// Get method returns value in cache by key.
// The second param in return will indicate if value by key exists or not.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	// Look up for the value by key
	v, isExist := c.data[key]
	if isExist && !v.isExpired() {
		// If value exists and not expired return the value
		return v.value, isExist
	}

	return nil, false
}

// Remove method removes the value in cache by key.
func (c *Cache) Remove(key string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	delete(c.data, key)
}

// Keys returns a list of all keys in cache.
func (c *Cache) Keys() []string {
	c.mux.RLock()
	defer c.mux.RUnlock()

	keys := make([]string, 0, len(c.data))
	for k, v := range c.data {
		if !v.isExpired() {
			keys = append(keys, k)
		}
	}

	return keys
}

// RPush method adds element to the list in cache.
// If key does not exist, a new key holding a list is created.
// TTL param could be omitted if it's adding to the existing list.
// If given TTL <=0 then the key will never be expired.
func (c *Cache) RPush(key string, value interface{}, ttl time.Duration) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	v, isExist := c.data[key]
	if !isExist || v.isExpired() {
		// Add new entity with list value
		e := entity{
			expiredAfter: validateExpiredAfter(ttl),
		}
		list := make([]interface{}, 0)
		list = append(list, value)
		e.value = list
		c.data[key] = e

		return nil
	}

	// Check if found value it's a slice
	sl, ok := v.value.([]interface{})
	if !ok {
		return ErrWrongTypeLPush
	}

	// Add new item to the slice and update entity in cache
	sl = append(sl, value)
	v.value = sl
	c.data[key] = v

	return nil
}

// LIndex method returns the element at the index in the list.
// The index is zero-based, 0 means the first element of the list.
// When the value at key is not a list, an error is returned.
// When index is not exist in the list - nil value is returned.
func (c *Cache) LIndex(key string, index int) (interface{}, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	v, isExist := c.data[key]
	if isExist && !v.isExpired() {
		// Check if type is slice
		sl, ok := v.value.([]interface{})
		if !ok {
			return nil, ErrWrongTypeIndex
		}

		// Check if index is exist and return nil value if it's not
		if len(sl)-1 < index {
			return nil, nil
		}

		return sl[index], nil
	}

	return nil, ErrNotFound
}

// HSet method sets value in the hash stored at key to value.
// If key does not exist, a new key holding a hash is created.
// If field already exists in the hash, it is overwritten.
// TTL param could be omitted if it's adding to the existing hash map.
func (c *Cache) HSet(key string, value map[string]interface{}, ttl time.Duration) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	v, isExist := c.data[key]
	if !isExist || v.isExpired() {
		// Add new entity with hm value
		e := entity{
			value:        value,
			expiredAfter: validateExpiredAfter(ttl),
		}
		c.data[key] = e

		return nil
	}

	// Check if found value it's a map
	hm, ok := v.value.(map[string]interface{})
	if !ok {
		return ErrWrongTypeHSet
	}

	// Avoid cases when hm has been initialized as nil
	if hm == nil {
		hm = make(map[string]interface{})
	}

	// Set new values to hash map
	for k, v := range value {
		hm[k] = v
	}

	// Update entity in cache
	v.value = hm
	c.data[key] = v

	return nil
}

// HGet method returns the value associated with field in the hash stored at key.
// When the value at key is not a hash map, an error is returned.
// When key in hash map value is not exist - nil value is returned.
func (c *Cache) HGet(key, hkey string) (interface{}, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	v, isExist := c.data[key]
	if isExist && !v.isExpired() {
		// Check if type is map
		hm, ok := v.value.(map[string]interface{})
		if !ok {
			return nil, ErrWrongTypeHGet
		}

		return hm[hkey], nil
	}

	return nil, ErrNotFound
}

func validateExpiredAfter(ttl time.Duration) int64 {
	var expiredAfter int64

	// Set time when cache is expired
	if ttl > 0 {
		expiredAfter = time.Now().UTC().Add(ttl).UnixNano()
	}

	return expiredAfter
}
