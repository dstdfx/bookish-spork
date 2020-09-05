package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/dstdfx/bookish-spork/httpclient/testutils"
	"github.com/stretchr/testify/require"
)

const (
	testKey               = "test-key"
	testHKey              = "test-hkey"
	testIndex             = 1
	testGetResponseRaw    = `{"value": "test-value"}`
	testSetRawRequest     = `{"key": "test-key", "value": "test-value", "ttl": 10}`
	testKeysResponseRaw   = `{"keys": ["test-key0", "test-key1", "test-key2"]}`
	testRPushRawRequest   = `{"key": "test-key", "value": "test-value", "ttl": 10}`
	testLIndexRawResponse = `{"value": "test-value"}`
	testHSetRawRequest    = `{"key": "test-key", "value": {"key0": "value0"}, "ttl": 10}`
	testHGetRawResponse   = `{"value": "hvalue"}`
)

var (
	expectedGet        interface{} = "test-value"
	expectedKeys                   = []string{"test-key0", "test-key1", "test-key2"}
	expectedIndexValue interface{} = "test-value"
	expectedHValue     interface{} = "hvalue"
)

func TestGet(t *testing.T) {
	endpointCalled := false
	testEnv := testutils.SetupTestEnv()
	defer testEnv.TearDownTestEnv()

	testutils.HandleReqWithoutBody(t, &testutils.HandleReqOpts{
		Mux:         testEnv.Mux,
		URL:         fmt.Sprintf("/v1/get/%s", testKey),
		RawResponse: testGetResponseRaw,
		Method:      http.MethodGet,
		Status:      http.StatusOK,
		CallFlag:    &endpointCalled,
	})

	ctx := context.Background()
	testClient := NewClient(testEnv.Server.URL + "/v1")

	actual, httpResponse, err := testClient.Get(ctx, testKey)
	require.NoError(t, err)
	require.True(t, endpointCalled)
	require.NotNil(t, httpResponse)
	require.Equal(t, http.StatusOK, httpResponse.StatusCode)
	require.Equal(t, expectedGet, actual)
}

func TestSet(t *testing.T) {
	endpointCalled := false
	testEnv := testutils.SetupTestEnv()
	defer testEnv.TearDownTestEnv()

	testutils.HandleReqWithoutBody(t, &testutils.HandleReqOpts{
		Mux:        testEnv.Mux,
		URL:        "/v1/set",
		RawRequest: testSetRawRequest,
		Method:     http.MethodPost,
		Status:     http.StatusOK,
		CallFlag:   &endpointCalled,
	})

	ctx := context.Background()
	testClient := NewClient(testEnv.Server.URL + "/v1")

	httpResponse, err := testClient.Set(ctx, SetBody{
		Key:   "test-key",
		Value: "test-value",
		TTL:   10,
	})
	require.NoError(t, err)
	require.True(t, endpointCalled)
	require.NotNil(t, httpResponse)
	require.Equal(t, http.StatusOK, httpResponse.StatusCode)
}

func TestKeys(t *testing.T) {
	endpointCalled := false
	testEnv := testutils.SetupTestEnv()
	defer testEnv.TearDownTestEnv()

	testutils.HandleReqWithoutBody(t, &testutils.HandleReqOpts{
		Mux:         testEnv.Mux,
		URL:         "/v1/keys",
		RawResponse: testKeysResponseRaw,
		Method:      http.MethodGet,
		Status:      http.StatusOK,
		CallFlag:    &endpointCalled,
	})

	ctx := context.Background()
	testClient := NewClient(testEnv.Server.URL + "/v1")

	actual, httpResponse, err := testClient.Keys(ctx)
	require.NoError(t, err)
	require.True(t, endpointCalled)
	require.NotNil(t, httpResponse)
	require.Equal(t, http.StatusOK, httpResponse.StatusCode)
	require.Equal(t, expectedKeys, actual)
}

func TestRemove(t *testing.T) {
	endpointCalled := false
	testEnv := testutils.SetupTestEnv()
	defer testEnv.TearDownTestEnv()

	testutils.HandleReqWithoutBody(t, &testutils.HandleReqOpts{
		Mux:      testEnv.Mux,
		URL:      fmt.Sprintf("/v1/remove/%s", testKey),
		Method:   http.MethodDelete,
		Status:   http.StatusNoContent,
		CallFlag: &endpointCalled,
	})

	ctx := context.Background()
	testClient := NewClient(testEnv.Server.URL + "/v1")

	httpResponse, err := testClient.Remove(ctx, testKey)
	require.NoError(t, err)
	require.True(t, endpointCalled)
	require.NotNil(t, httpResponse)
	require.Equal(t, http.StatusNoContent, httpResponse.StatusCode)
}

func TestRPush(t *testing.T) {
	endpointCalled := false
	testEnv := testutils.SetupTestEnv()
	defer testEnv.TearDownTestEnv()

	testutils.HandleReqWithoutBody(t, &testutils.HandleReqOpts{
		Mux:        testEnv.Mux,
		URL:        "/v1/rpush",
		RawRequest: testRPushRawRequest,
		Method:     http.MethodPost,
		Status:     http.StatusOK,
		CallFlag:   &endpointCalled,
	})

	ctx := context.Background()
	testClient := NewClient(testEnv.Server.URL + "/v1")

	httpResponse, err := testClient.RPush(ctx, RPushBody{
		Key:   "test-key",
		Value: "test-value",
		TTL:   10,
	})
	require.NoError(t, err)
	require.True(t, endpointCalled)
	require.NotNil(t, httpResponse)
	require.Equal(t, http.StatusOK, httpResponse.StatusCode)
}

func TestLIndex(t *testing.T) {
	endpointCalled := false
	testEnv := testutils.SetupTestEnv()
	defer testEnv.TearDownTestEnv()

	testutils.HandleReqWithoutBody(t, &testutils.HandleReqOpts{
		Mux:         testEnv.Mux,
		URL:         fmt.Sprintf("/v1/lindex/%s/%d", testKey, testIndex),
		RawResponse: testLIndexRawResponse,
		Method:      http.MethodGet,
		Status:      http.StatusOK,
		CallFlag:    &endpointCalled,
	})

	ctx := context.Background()
	testClient := NewClient(testEnv.Server.URL + "/v1")

	actual, httpResponse, err := testClient.LIndex(ctx, testKey, testIndex)
	require.NoError(t, err)
	require.True(t, endpointCalled)
	require.NotNil(t, httpResponse)
	require.Equal(t, http.StatusOK, httpResponse.StatusCode)
	require.Equal(t, expectedIndexValue, actual)
}

func TestHSet(t *testing.T) {
	endpointCalled := false
	testEnv := testutils.SetupTestEnv()
	defer testEnv.TearDownTestEnv()

	testutils.HandleReqWithoutBody(t, &testutils.HandleReqOpts{
		Mux:        testEnv.Mux,
		URL:        "/v1/hset",
		RawRequest: testHSetRawRequest,
		Method:     http.MethodPost,
		Status:     http.StatusOK,
		CallFlag:   &endpointCalled,
	})

	ctx := context.Background()
	testClient := NewClient(testEnv.Server.URL + "/v1")

	httpResponse, err := testClient.HSet(ctx, HSetBody{
		Key:   "test-key",
		Value: map[string]interface{}{"key0": "value0"},
		TTL:   10,
	})
	require.NoError(t, err)
	require.True(t, endpointCalled)
	require.NotNil(t, httpResponse)
	require.Equal(t, http.StatusOK, httpResponse.StatusCode)
}

func TestHGet(t *testing.T) {
	endpointCalled := false
	testEnv := testutils.SetupTestEnv()
	defer testEnv.TearDownTestEnv()

	testutils.HandleReqWithoutBody(t, &testutils.HandleReqOpts{
		Mux:         testEnv.Mux,
		URL:         fmt.Sprintf("/v1/hget/%s/%s", testKey, testHKey),
		RawResponse: testHGetRawResponse,
		Method:      http.MethodGet,
		Status:      http.StatusOK,
		CallFlag:    &endpointCalled,
	})

	ctx := context.Background()
	testClient := NewClient(testEnv.Server.URL + "/v1")

	actual, httpResponse, err := testClient.HGet(ctx, testKey, testHKey)
	require.NoError(t, err)
	require.True(t, endpointCalled)
	require.NotNil(t, httpResponse)
	require.Equal(t, http.StatusOK, httpResponse.StatusCode)
	require.Equal(t, expectedHValue, actual)
}
