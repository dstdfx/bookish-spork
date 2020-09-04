package qqcache

import (
	"sync"
	"time"
)

const defaultEvictionInterval = 60 * time.Second

// TODO: implement cache cleaner
// TODO: add persistent keys

// Cache represents cache container.
type Cache struct {
	mux              sync.RWMutex
	data             map[string]entity
	evictionInterval time.Duration
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
	}

	return c
}

// entity represents an object that stores by key in cache.
type entity struct {
	value        interface{}
	expiredAfter int64
}

// Set method sets value to cache by key with specific TTL.
// If given TTL <=0 then the key will never be expired.
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()

	var expiredAfter int64

	// Set time when cache is expired
	if ttl > 0 {
		expiredAfter = time.Now().UTC().Add(ttl).UnixNano()
	}

	c.data[key] = entity{
		value:        value,
		expiredAfter: expiredAfter,
	}
}

// Get method returns value in cache by key.
// The second param in return will indicate if value by key exists or not.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	// Look up for the value by key
	v, isExist := c.data[key]
	if isExist && time.Now().UTC().UnixNano() < v.expiredAfter {
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
		if time.Now().UTC().UnixNano() < v.expiredAfter {
			keys = append(keys, k)
		}
	}

	return keys
}
