package qqcache

import (
	"errors"
	"sync"
	"time"
)

const defaultEvictionInterval = 60 * time.Second

var (
	ErrWrongTypeLPush = errors.New("wrong type of the value to push list value")
)

// Cache represents cache container.
type Cache struct {
	mux              sync.RWMutex
	data             map[string]entity
	evictionInterval time.Duration
	stopCleaner      chan struct{}
}

// New returns new instance of Cache.
// If eviction interval is equal or less that 0 - default eviction will be used.
func New(evictionInterval time.Duration) *Cache {
	if evictionInterval <= 0 {
		evictionInterval = defaultEvictionInterval
	}

	c := &Cache{
		mux:              sync.RWMutex{},
		data:             make(map[string]entity),
		evictionInterval: evictionInterval,
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

func validateExpiredAfter(ttl time.Duration) int64 {
	var expiredAfter int64

	// Set time when cache is expired
	if ttl > 0 {
		expiredAfter = time.Now().UTC().Add(ttl).UnixNano()
	}

	return expiredAfter
}
