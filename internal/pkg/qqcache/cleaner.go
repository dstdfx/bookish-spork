package qqcache

import "time"

// cacheCleaner runs cleaner that will delete expired keys within each
// eviction interval.
func (c *Cache) cacheCleaner() {
	t := time.NewTicker(c.evictionInterval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			// Stop cleaner if cache data is not initialized
			if c.data == nil {
				return
			}

			// Get a slice of expired keys
			expiredKeys := c.getExpiredKeys()
			if len(expiredKeys) == 0 {
				// Skip if there's no keys to delete
				continue
			}

			// Delete expired keys from the cache
			c.deleteExpiredKeys(expiredKeys)
		case <-c.stopCleaner:
			return
		}
	}
}

// getExpiredKeys method returns all expired keys in cache.
func (c *Cache) getExpiredKeys() []string {
	c.mux.RLock()
	defer c.mux.RUnlock()

	expiredKeys := make([]string, 0)
	for k, v := range c.data {
		if v.isExpired() {
			expiredKeys = append(expiredKeys, k)
		}
	}

	return expiredKeys
}

// deleteExpiredKeys method deletes given keys.
func (c *Cache) deleteExpiredKeys(expiredKeys []string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	for _, k := range expiredKeys {
		delete(c.data, k)
	}
}
