package qqcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	testKey   = "test-key"
	testValue = "test-value"
)

func TestNew_DefaultEviction(t *testing.T) {
	c := New(0)
	require.NotEmpty(t, c)
	require.Equal(t, defaultEvictionInterval, c.evictionInterval)
}

func TestCache_Get_Set(t *testing.T) {
	c := New(10 * time.Second)

	// Set test value to the cache
	c.Set(testKey, testValue, time.Second)

	// Get value from the cache
	got, ok := c.Get(testKey)

	// Check test value exists by the key and it's the same string
	require.True(t, ok)
	gotString, ok := got.(string)
	require.True(t, ok)
	require.Equal(t, testValue, gotString)

	// Wait for 3 seconds to let the key get expired
	<-time.After(time.Second * 3)

	// Try to get expired value by the key from the cache
	got, ok = c.Get(testKey)
	require.False(t, ok)
	require.Nil(t, got)
}

func TestCache_Get_Set_PersistentKey(t *testing.T) {
	c := New(10 * time.Second)

	// Set test value to the cache with 0 ttl
	c.Set(testKey, testValue, 0)

	// Get value from the cache
	got, ok := c.Get(testKey)

	// Check test value exists by the key and it's the same string
	require.True(t, ok)
	gotString, ok := got.(string)
	require.True(t, ok)
	require.Equal(t, testValue, gotString)
}

func TestCache_SetNil(t *testing.T) {
	c := New(10 * time.Second)

	// Set nil value to the cache
	c.Set(testKey, nil, time.Second)

	// Get nil value from the cache
	got, ok := c.Get(testKey)
	require.True(t, ok)
	require.Nil(t, got)
}

func TestCache_Remove(t *testing.T) {
	c := New(10 * time.Second)

	// Set test value to the cache
	c.Set(testKey, testValue, time.Second)

	// Remove value from the cache
	c.Remove(testKey)

	// Try to get deleted value by the key from the cache
	got, ok := c.Get(testKey)
	require.False(t, ok)
	require.Nil(t, got)
}

func TestCache_Remove_NoKey(t *testing.T) {
	c := New(10 * time.Second)
	testKey := "test-key"

	// Remove value from the cache
	c.Remove(testKey)
}

func TestCache_Keys(t *testing.T) {
	c := New(10 * time.Second)
	keysToSet := map[string]string{
		"test-key-0": "test-value-0",
		"test-key-1": "test-value-1",
		"test-key-2": "test-value-2",
	}
	expected := []string{"test-key-0", "test-key-1", "test-key-2"}

	for k, v := range keysToSet {
		c.Set(k, v, 5*time.Second)
	}

	require.ElementsMatch(t, expected, c.Keys())
}
