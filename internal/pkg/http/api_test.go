package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dstdfx/bookish-spork/internal/pkg/backend"
	"github.com/dstdfx/bookish-spork/internal/pkg/config"
	v1 "github.com/dstdfx/bookish-spork/internal/pkg/http/v1"
	"github.com/dstdfx/bookish-spork/internal/pkg/log"
	"github.com/dstdfx/bookish-spork/internal/pkg/qqcache"
	"github.com/dstdfx/bookish-spork/internal/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

const (
	testKey   = "test-key"
	testValue = "test-value"
)

// Tests for POST /v1/set

func TestSet_OK(t *testing.T) {
	// Check acceptance test flag
	if !testutils.IsAccTestEnabled(t) {
		return
	}

	// Init global app configuration
	testutils.InitTestConfig()

	// Initialize logger
	logger, err := log.InitLogger(log.InitLoggerOpts{
		Debug:     config.Config.Log.Debug,
		UseStdout: config.Config.Log.UseStdout,
		File:      config.Config.Log.File,
	})
	assert.NoError(t, err)

	// Prepare backend
	b := backend.New(logger)
	defer b.Shutdown()
	assert.NotEmpty(t, b)

	setBody := &v1.SetRequestBody{
		Key:   testKey,
		Value: testValue,
		TTL:   10,
	}
	reqBody, err := json.Marshal(setBody)
	assert.NoError(t, err)

	// Setup handlers
	router := InitAPIRouter(b)

	// Test a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodPost, "/v1/set", bytes.NewReader(reqBody))
	assert.NoError(t, err)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSet_BadRequest(t *testing.T) {
	// Check acceptance test flag
	if !testutils.IsAccTestEnabled(t) {
		return
	}

	// Init global app configuration
	testutils.InitTestConfig()

	// Initialize logger
	logger, err := log.InitLogger(log.InitLoggerOpts{
		Debug:     config.Config.Log.Debug,
		UseStdout: config.Config.Log.UseStdout,
		File:      config.Config.Log.File,
	})
	assert.NoError(t, err)

	// Prepare backend
	b := backend.New(logger)
	defer b.Shutdown()
	assert.NotEmpty(t, b)

	setBody := &v1.SetRequestBody{
		Key: testKey,
	}
	reqBody, err := json.Marshal(setBody)
	assert.NoError(t, err)

	// Setup handlers
	router := InitAPIRouter(b)

	// Test a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodPost, "/v1/set", bytes.NewReader(reqBody))
	assert.NoError(t, err)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Tests for GET /v1/get/<key>

func TestGet_OK(t *testing.T) {
	// Check acceptance test flag
	if !testutils.IsAccTestEnabled(t) {
		return
	}

	// Init global app configuration
	testutils.InitTestConfig()

	// Initialize logger
	logger, err := log.InitLogger(log.InitLoggerOpts{
		Debug:     config.Config.Log.Debug,
		UseStdout: config.Config.Log.UseStdout,
		File:      config.Config.Log.File,
	})
	assert.NoError(t, err)

	// Prepare backend
	b := backend.New(logger)
	assert.NotEmpty(t, b)

	// Set test value to cache
	b.Cache.Set(testKey, testValue, 0)

	// Setup handlers
	router := InitAPIRouter(b)

	// Test a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/get/%s", testKey), nil)
	assert.NoError(t, err)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t,
		testutils.RespToJSON(t,
			map[string]string{"value": testValue},
		), w.Body.String())
}

func TestGet_NotFound(t *testing.T) {
	// Check acceptance test flag
	if !testutils.IsAccTestEnabled(t) {
		return
	}

	// Init global app configuration
	testutils.InitTestConfig()

	// Initialize logger
	logger, err := log.InitLogger(log.InitLoggerOpts{
		Debug:     config.Config.Log.Debug,
		UseStdout: config.Config.Log.UseStdout,
		File:      config.Config.Log.File,
	})
	assert.NoError(t, err)

	// Prepare backend
	b := backend.New(logger)
	defer b.Shutdown()
	assert.NotEmpty(t, b)

	// Setup handlers
	router := InitAPIRouter(b)

	// Test a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/get/%s", testKey), nil)
	assert.NoError(t, err)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// Tests for GET /v1/keys

func TestKeys_OK(t *testing.T) {
	// Check acceptance test flag
	if !testutils.IsAccTestEnabled(t) {
		return
	}

	// Init global app configuration
	testutils.InitTestConfig()

	// Initialize logger
	logger, err := log.InitLogger(log.InitLoggerOpts{
		Debug:     config.Config.Log.Debug,
		UseStdout: config.Config.Log.UseStdout,
		File:      config.Config.Log.File,
	})
	assert.NoError(t, err)

	// Prepare backend
	b := backend.New(logger)
	defer b.Shutdown()
	assert.NotEmpty(t, b)

	// Set test value to cache
	b.Cache.Set(testKey, testValue, 0)

	// Setup handlers
	router := InitAPIRouter(b)

	// Test a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/v1/keys", nil)
	assert.NoError(t, err)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t,
		testutils.RespToJSON(t,
			map[string][]string{"keys": {testKey}},
		), w.Body.String())
}

// Tests for DELETE /v1/remove/<key>

func TestRemove_OK(t *testing.T) {
	// Check acceptance test flag
	if !testutils.IsAccTestEnabled(t) {
		return
	}

	// Init global app configuration
	testutils.InitTestConfig()

	// Initialize logger
	logger, err := log.InitLogger(log.InitLoggerOpts{
		Debug:     config.Config.Log.Debug,
		UseStdout: config.Config.Log.UseStdout,
		File:      config.Config.Log.File,
	})
	assert.NoError(t, err)

	// Prepare backend
	b := backend.New(logger)
	defer b.Shutdown()
	assert.NotEmpty(t, b)

	// Set test value to cache
	b.Cache.Set(testKey, testValue, 0)

	// Setup handlers
	router := InitAPIRouter(b)

	// Test a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/remove/%s", testKey), nil)
	assert.NoError(t, err)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// Tests for POST /v1/rpush

func TestRPush_OK(t *testing.T) {
	// Check acceptance test flag
	if !testutils.IsAccTestEnabled(t) {
		return
	}

	// Init global app configuration
	testutils.InitTestConfig()

	// Initialize logger
	logger, err := log.InitLogger(log.InitLoggerOpts{
		Debug:     config.Config.Log.Debug,
		UseStdout: config.Config.Log.UseStdout,
		File:      config.Config.Log.File,
	})
	assert.NoError(t, err)

	// Prepare backend
	b := backend.New(logger)
	defer b.Shutdown()
	assert.NotEmpty(t, b)

	rpushBody := &v1.RPushRequestBody{
		Key:   testKey,
		Value: testValue,
		TTL:   10,
	}
	reqBody, err := json.Marshal(rpushBody)
	assert.NoError(t, err)

	// Setup handlers
	router := InitAPIRouter(b)

	// Test a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodPost, "/v1/rpush", bytes.NewReader(reqBody))
	assert.NoError(t, err)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRPush_AppendValue(t *testing.T) {
	// Check acceptance test flag
	if !testutils.IsAccTestEnabled(t) {
		return
	}

	// Init global app configuration
	testutils.InitTestConfig()

	// Initialize logger
	logger, err := log.InitLogger(log.InitLoggerOpts{
		Debug:     config.Config.Log.Debug,
		UseStdout: config.Config.Log.UseStdout,
		File:      config.Config.Log.File,
	})
	assert.NoError(t, err)

	// Prepare backend
	b := backend.New(logger)
	defer b.Shutdown()
	assert.NotEmpty(t, b)

	assert.NoError(t, b.Cache.RPush(testKey, testValue, 0))

	rpushBody := &v1.RPushRequestBody{
		Key:   testKey,
		Value: testValue,
		TTL:   10,
	}
	reqBody, err := json.Marshal(rpushBody)
	assert.NoError(t, err)

	// Setup handlers
	router := InitAPIRouter(b)

	// Test a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodPost, "/v1/rpush", bytes.NewReader(reqBody))
	assert.NoError(t, err)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRPush_BadRequest(t *testing.T) {
	// Check acceptance test flag
	if !testutils.IsAccTestEnabled(t) {
		return
	}

	// Init global app configuration
	testutils.InitTestConfig()

	// Initialize logger
	logger, err := log.InitLogger(log.InitLoggerOpts{
		Debug:     config.Config.Log.Debug,
		UseStdout: config.Config.Log.UseStdout,
		File:      config.Config.Log.File,
	})
	assert.NoError(t, err)

	// Prepare backend
	b := backend.New(logger)
	defer b.Shutdown()
	assert.NotEmpty(t, b)

	// Test non-list value
	b.Cache.Set(testKey, 1, 0)

	rpushBody := &v1.RPushRequestBody{
		Key:   testKey,
		Value: testValue,
	}
	reqBody, err := json.Marshal(rpushBody)
	assert.NoError(t, err)

	// Setup handlers
	router := InitAPIRouter(b)

	// Test a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodPost, "/v1/rpush", bytes.NewReader(reqBody))
	assert.NoError(t, err)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Tests for GET /v1/lindex/<key>/<index>

func TestLindex_OK(t *testing.T) {
	// Check acceptance test flag
	if !testutils.IsAccTestEnabled(t) {
		return
	}

	// Init global app configuration
	testutils.InitTestConfig()

	// Initialize logger
	logger, err := log.InitLogger(log.InitLoggerOpts{
		Debug:     config.Config.Log.Debug,
		UseStdout: config.Config.Log.UseStdout,
		File:      config.Config.Log.File,
	})
	assert.NoError(t, err)

	// Prepare backend
	b := backend.New(logger)
	assert.NotEmpty(t, b)

	// Set test value to cache
	assert.NoError(t, b.Cache.RPush(testKey, testValue, 0))

	// Setup handlers
	router := InitAPIRouter(b)

	// Test a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/lindex/%s/%d", testKey, 0), nil)
	assert.NoError(t, err)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t,
		testutils.RespToJSON(t,
			map[string]string{"value": testValue},
		), w.Body.String())
}

func TestLindex_BadRequest(t *testing.T) {
	// Check acceptance test flag
	if !testutils.IsAccTestEnabled(t) {
		return
	}

	// Init global app configuration
	testutils.InitTestConfig()

	// Initialize logger
	logger, err := log.InitLogger(log.InitLoggerOpts{
		Debug:     config.Config.Log.Debug,
		UseStdout: config.Config.Log.UseStdout,
		File:      config.Config.Log.File,
	})
	assert.NoError(t, err)

	// Prepare backend
	b := backend.New(logger)
	assert.NotEmpty(t, b)

	// Set test value to cache
	b.Cache.Set(testKey, testValue, 0)

	// Setup handlers
	router := InitAPIRouter(b)

	// Test a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/lindex/%s/%d", testKey, 1), nil)
	assert.NoError(t, err)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t,
		testutils.RespToJSON(t,
			map[string]string{"error": qqcache.ErrWrongTypeIndex.Error()},
		), w.Body.String())
}

func TestLindex_BadIndex(t *testing.T) {
	// Check acceptance test flag
	if !testutils.IsAccTestEnabled(t) {
		return
	}

	// Init global app configuration
	testutils.InitTestConfig()

	// Initialize logger
	logger, err := log.InitLogger(log.InitLoggerOpts{
		Debug:     config.Config.Log.Debug,
		UseStdout: config.Config.Log.UseStdout,
		File:      config.Config.Log.File,
	})
	assert.NoError(t, err)

	// Prepare backend
	b := backend.New(logger)
	assert.NotEmpty(t, b)

	// Setup handlers
	router := InitAPIRouter(b)

	// Test a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/lindex/%s/%d", testKey, -1), nil)
	assert.NoError(t, err)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t,
		testutils.RespToJSON(t,
			map[string]string{"error": "index can't be negative"},
		), w.Body.String())
}

// Tests for POST /v1/hset

func TestHSet_OK(t *testing.T) {
	// Check acceptance test flag
	if !testutils.IsAccTestEnabled(t) {
		return
	}

	// Init global app configuration
	testutils.InitTestConfig()

	// Initialize logger
	logger, err := log.InitLogger(log.InitLoggerOpts{
		Debug:     config.Config.Log.Debug,
		UseStdout: config.Config.Log.UseStdout,
		File:      config.Config.Log.File,
	})
	assert.NoError(t, err)

	// Prepare backend
	b := backend.New(logger)
	defer b.Shutdown()
	assert.NotEmpty(t, b)

	rpushBody := &v1.HSetRequestBody{
		Key:   testKey,
		Value: map[string]interface{}{"some-key": "some-value"},
		TTL:   10,
	}
	reqBody, err := json.Marshal(rpushBody)
	assert.NoError(t, err)

	// Setup handlers
	router := InitAPIRouter(b)

	// Test a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodPost, "/v1/hset", bytes.NewReader(reqBody))
	assert.NoError(t, err)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHSet_BadRequest(t *testing.T) {
	// Check acceptance test flag
	if !testutils.IsAccTestEnabled(t) {
		return
	}

	// Init global app configuration
	testutils.InitTestConfig()

	// Initialize logger
	logger, err := log.InitLogger(log.InitLoggerOpts{
		Debug:     config.Config.Log.Debug,
		UseStdout: config.Config.Log.UseStdout,
		File:      config.Config.Log.File,
	})
	assert.NoError(t, err)

	// Prepare backend
	b := backend.New(logger)
	defer b.Shutdown()
	assert.NotEmpty(t, b)

	rpushBody := &v1.HSetRequestBody{
		Key: testKey,
	}
	reqBody, err := json.Marshal(rpushBody)
	assert.NoError(t, err)

	// Setup handlers
	router := InitAPIRouter(b)

	// Test a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodPost, "/v1/hset", bytes.NewReader(reqBody))
	assert.NoError(t, err)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHSet_AppendKeys(t *testing.T) {
	// Check acceptance test flag
	if !testutils.IsAccTestEnabled(t) {
		return
	}

	// Init global app configuration
	testutils.InitTestConfig()

	// Initialize logger
	logger, err := log.InitLogger(log.InitLoggerOpts{
		Debug:     config.Config.Log.Debug,
		UseStdout: config.Config.Log.UseStdout,
		File:      config.Config.Log.File,
	})
	assert.NoError(t, err)

	// Prepare backend
	b := backend.New(logger)
	defer b.Shutdown()
	assert.NotEmpty(t, b)

	assert.NoError(t, b.Cache.HSet(testKey, map[string]interface{}{"some-key-0": "some-value-0"}, 0))

	rpushBody := &v1.HSetRequestBody{
		Key:   testKey,
		Value: map[string]interface{}{"some-key-1": "some-value-1"},
	}
	reqBody, err := json.Marshal(rpushBody)
	assert.NoError(t, err)

	// Setup handlers
	router := InitAPIRouter(b)

	// Test a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodPost, "/v1/hset", bytes.NewReader(reqBody))
	assert.NoError(t, err)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Tests for GET /v1/hget/<key>/<hkey>
