package qqcache

import (
	"errors"
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
	defer c.Shutdown()
	require.NotEmpty(t, c)
	require.Equal(t, defaultEvictionInterval, c.evictionInterval)
}

func TestCache_Get_Set(t *testing.T) {
	c := New(10 * time.Second)
	defer c.Shutdown()

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
	defer c.Shutdown()

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
	defer c.Shutdown()

	// Set nil value to the cache
	c.Set(testKey, nil, time.Second)

	// Get nil value from the cache
	got, ok := c.Get(testKey)
	require.True(t, ok)
	require.Nil(t, got)
}

func TestCache_Remove(t *testing.T) {
	c := New(10 * time.Second)
	defer c.Shutdown()

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
	defer c.Shutdown()

	// Remove value from the cache
	c.Remove(testKey)
}

func TestCache_Keys(t *testing.T) {
	c := New(10 * time.Second)
	defer c.Shutdown()
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

func TestCache_RPush(t *testing.T) {
	c := New(10 * time.Second)
	expected := []int{1, 2, 3, 4, 5}

	// Add list and values
	require.NoError(t, c.RPush(testKey, 1, 0))
	require.NoError(t, c.RPush(testKey, 2, 0))
	require.NoError(t, c.RPush(testKey, 3, 0))
	require.NoError(t, c.RPush(testKey, 4, 0))
	require.NoError(t, c.RPush(testKey, 5, 0))

	// Check list has been added
	v, ok := c.Get(testKey)
	require.True(t, ok)
	require.NotEmpty(t, v)

	// Check values
	sl, ok := v.([]interface{})
	require.True(t, ok)

	got := make([]int, len(sl))
	for i := range sl {
		got[i] = sl[i].(int)
	}
	require.EqualValues(t, expected, got)
}

func TestCache_RPush_Expired(t *testing.T) {
	c := New(10 * time.Second)
	expected := []int{3, 5}

	// Add some values to list
	require.NoError(t, c.RPush(testKey, 1, time.Second))
	require.NoError(t, c.RPush(testKey, 2, 0))
	require.NoError(t, c.RPush(testKey, 3, 0))
	require.NoError(t, c.RPush(testKey, 4, 0))
	require.NoError(t, c.RPush(testKey, 5, 0))

	// Wait til the key is expired
	<-time.After(3 * time.Second)

	// Check no value by the key
	v, ok := c.Get(testKey)
	require.False(t, ok)
	require.Empty(t, v)

	// Add new list
	require.NoError(t, c.RPush(testKey, 3, 0))
	require.NoError(t, c.RPush(testKey, 5, 0))

	// Check list has been added
	v, ok = c.Get(testKey)
	require.True(t, ok)
	require.NotEmpty(t, v)

	// Check values of the list
	sl, ok := v.([]interface{})
	require.True(t, ok)

	got := make([]int, len(sl))
	for i := range sl {
		got[i] = sl[i].(int)
	}
	require.EqualValues(t, expected, got)
}

func TestCache_RPush_WrongValueType(t *testing.T) {
	c := New(10 * time.Second)

	// Add string value
	c.Set(testKey, testValue, 0)

	// Try to add element to the list
	err := c.RPush(testKey, 1, 0)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrWrongTypeLPush))
}
