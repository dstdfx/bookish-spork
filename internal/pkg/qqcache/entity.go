package qqcache

import "time"

// entity represents an object that stores by key in cache.
type entity struct {
	value        interface{}
	expiredAfter int64
}

// isExpired method returns true if the value is expired.
func (e entity) isExpired() bool {
	// Check if value is set to be persistent
	if e.expiredAfter <= 0 {
		return false
	}

	return time.Now().UTC().UnixNano() > e.expiredAfter
}
