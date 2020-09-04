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

func getCommonCacheOpts() Opts {
	return Opts{EvictionInterval: 10 * time.Second}
}

func TestNew_DefaultEviction(t *testing.T) {
	c := New(Opts{EvictionInterval: 0})
	defer c.Shutdown()
	require.NotEmpty(t, c)
	require.Equal(t, defaultEvictionInterval, c.evictionInterval)
}

func TestCache_Get_Set(t *testing.T) {
	c := New(getCommonCacheOpts())
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
	c := New(getCommonCacheOpts())
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
	c := New(getCommonCacheOpts())
	defer c.Shutdown()

	// Set nil value to the cache
	c.Set(testKey, nil, time.Second)

	// Get nil value from the cache
	got, ok := c.Get(testKey)
	require.True(t, ok)
	require.Nil(t, got)
}

func TestCache_Remove(t *testing.T) {
	c := New(getCommonCacheOpts())
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
	c := New(getCommonCacheOpts())
	defer c.Shutdown()

	// Remove value from the cache
	c.Remove(testKey)
}

func TestCache_Keys(t *testing.T) {
	c := New(getCommonCacheOpts())
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
	c := New(getCommonCacheOpts())
	defer c.Shutdown()
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
	c := New(getCommonCacheOpts())
	defer c.Shutdown()
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
	c := New(getCommonCacheOpts())
	defer c.Shutdown()

	// Add string value
	c.Set(testKey, testValue, 0)

	// Try to add element to the list
	err := c.RPush(testKey, 1, 0)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrWrongTypeLPush))
}

func TestCache_RPush_Nil(t *testing.T) {
	c := New(getCommonCacheOpts())
	defer c.Shutdown()

	// Add nil element to the list
	err := c.RPush(testKey, nil, 0)
	require.NoError(t, err)

	err = c.RPush(testKey, 1, 0)
	require.NoError(t, err)

	got, err := c.LIndex(testKey, 0)
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestCache_LIndex(t *testing.T) {
	c := New(getCommonCacheOpts())
	defer c.Shutdown()

	// Add list and values
	require.NoError(t, c.RPush(testKey, 1, 0))
	require.NoError(t, c.RPush(testKey, 2, 0))
	require.NoError(t, c.RPush(testKey, 3, 0))
	require.NoError(t, c.RPush(testKey, 4, 0))
	require.NoError(t, c.RPush(testKey, 5, 0))

	// Get value at the index 1
	got, err := c.LIndex(testKey, 1)
	require.NoError(t, err)
	require.NotEmpty(t, got)
	require.Equal(t, 2, got.(int))

	// Get value at the index 4
	got, err = c.LIndex(testKey, 4)
	require.NoError(t, err)
	require.NotEmpty(t, got)
	require.Equal(t, 5, got.(int))

	// Get value at not-existing index in the list
	got, err = c.LIndex(testKey, 5)
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestCache_LIndex_WrongValueType(t *testing.T) {
	c := New(getCommonCacheOpts())
	defer c.Shutdown()

	// Add string value
	c.Set(testKey, testValue, 0)

	// Try to get index of the list
	v, err := c.LIndex(testKey, 1)
	require.Error(t, err)
	require.Nil(t, v)
	require.True(t, errors.Is(err, ErrWrongTypeIndex))
}

func TestCache_LIndex_NotFound(t *testing.T) {
	c := New(getCommonCacheOpts())
	defer c.Shutdown()

	// Try to get index of the list
	v, err := c.LIndex(testKey, 1)
	require.Error(t, err)
	require.Nil(t, v)
	require.True(t, errors.Is(err, ErrNotFound))
}

func TestCache_HSet(t *testing.T) {
	c := New(getCommonCacheOpts())
	defer c.Shutdown()

	hmValue := map[string]interface{}{
		"key0": "value0",
		"key1": "value1",
		"key2": "value2",
	}

	// Set hm value
	err := c.HSet(testKey, hmValue, 0)
	require.NoError(t, err)

	// Get hm value to check
	v, ok := c.Get(testKey)
	require.True(t, ok)
	require.NotNil(t, v)

	// Check values of the hm
	got, ok := v.(map[string]interface{})
	require.True(t, ok)
	require.EqualValues(t, hmValue, got)
}

func TestCache_HSet_OverwriteHKeys(t *testing.T) {
	c := New(getCommonCacheOpts())
	defer c.Shutdown()

	hmValue := map[string]interface{}{
		"key0": "value0",
		"key1": "value1",
		"key2": "value2",
	}

	// Set hm value
	err := c.HSet(testKey, hmValue, 0)
	require.NoError(t, err)

	// Overwrite hm key value
	err = c.HSet(testKey, map[string]interface{}{"key0": "new-value"}, 0)
	require.NoError(t, err)

	// Get hm value to check
	v, ok := c.Get(testKey)
	require.True(t, ok)
	require.NotNil(t, v)

	expected := hmValue
	expected["key0"] = "new-value"

	// Check values of the hm
	got, ok := v.(map[string]interface{})
	require.True(t, ok)
	require.EqualValues(t, expected, got)
}

func TestCache_HSet_WrongType(t *testing.T) {
	c := New(getCommonCacheOpts())
	defer c.Shutdown()

	c.Set(testKey, 0, 0)

	// Try to set hm key into wrong value type
	err := c.HSet(testKey, map[string]interface{}{"key0": "new-value"}, 0)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrWrongTypeHSet))
}

func TestCache_HSet_NilHM(t *testing.T) {
	c := New(getCommonCacheOpts())
	defer c.Shutdown()

	err := c.HSet(testKey, nil, 0)
	require.NoError(t, err)

	// Set hm value to nil hm
	err = c.HSet(testKey, map[string]interface{}{"key0": "new-value"}, 0)
	require.NoError(t, err)

	// Get hm value to check
	v, ok := c.Get(testKey)
	require.True(t, ok)
	require.NotNil(t, v)

	expected := map[string]interface{}{"key0": "new-value"}

	// Check values of the hm
	got, ok := v.(map[string]interface{})
	require.True(t, ok)
	require.EqualValues(t, expected, got)
}

func TestCache_HGet(t *testing.T) {
	c := New(getCommonCacheOpts())
	defer c.Shutdown()

	hmValue := map[string]interface{}{
		"key0": "value0",
		"key1": "value1",
		"key2": "value2",
	}

	// Set hm value
	err := c.HSet(testKey, hmValue, 0)
	require.NoError(t, err)

	got, err := c.HGet(testKey, "key1")
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "value1", got.(string))
}

func TestCache_HGet_NoKey(t *testing.T) {
	c := New(getCommonCacheOpts())
	defer c.Shutdown()

	got, err := c.HGet(testKey, "some")
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrNotFound))
	require.Nil(t, got)
}

func TestCache_HGet_NoHKey(t *testing.T) {
	c := New(getCommonCacheOpts())
	defer c.Shutdown()

	hmValue := map[string]interface{}{
		"key0": "value0",
		"key1": "value1",
		"key2": "value2",
	}

	// Set hm value
	err := c.HSet(testKey, hmValue, 0)
	require.NoError(t, err)

	got, err := c.HGet(testKey, "not-existing")
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestCache_HGet_WrongType(t *testing.T) {
	c := New(getCommonCacheOpts())
	defer c.Shutdown()

	// Set value
	c.Set(testKey, 0, 0)

	got, err := c.HGet(testKey, "key1")
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrWrongTypeHGet))
	require.Nil(t, got)
}
