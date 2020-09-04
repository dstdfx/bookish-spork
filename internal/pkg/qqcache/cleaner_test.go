package qqcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRunCleaner(t *testing.T) {
	c := New(time.Second)

	// Set key with low TTL
	c.Set(testKey, testValue, time.Second)

	// Wait til TTL is expired
	<-time.After(3 * time.Second)

	// Stop cache cleaner to prevent locking
	c.Shutdown()

	// Check that key has been deleted by cache cleaner
	_, ok := c.data[testKey]
	require.False(t, ok)
}
