package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/dstdfx/bookish-spork/httpclient/testutils"
	"github.com/stretchr/testify/require"
)

func TestDoGetRequest(t *testing.T) {
	testEnv := testutils.SetupTestEnv()
	defer testEnv.TearDownTestEnv()
	testEnv.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, "response")

		require.Equal(t, http.MethodGet, r.Method)
	})

	endpoint := testEnv.Server.URL + "/"
	client := &Client{
		HTTPClient: &http.Client{},
		Endpoint:   endpoint,
	}

	ctx := context.Background()
	response, err := client.doRequest(ctx, http.MethodGet, endpoint, nil)
	require.NoError(t, err)
	require.NotNil(t, response.Body)
	require.Equal(t, http.StatusOK, response.StatusCode)
}

func TestDoPostRequest(t *testing.T) {
	testEnv := testutils.SetupTestEnv()
	defer testEnv.TearDownTestEnv()
	testEnv.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, "response")

		require.Equal(t, http.MethodPost, r.Method)

		_, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
	})

	endpoint := testEnv.Server.URL + "/"
	client := &Client{
		HTTPClient: &http.Client{},
		Endpoint:   endpoint,
	}

	requestBody, err := json.Marshal(&struct {
		ID string `json:"id"`
	}{
		ID: "uuid",
	})
	require.NoError(t, err)

	ctx := context.Background()
	response, err := client.doRequest(ctx, http.MethodPost, endpoint, bytes.NewReader(requestBody))
	require.NoError(t, err)
	require.NotNil(t, response.Body)
	require.Equal(t, http.StatusOK, response.StatusCode)
}

func TestDoErrNotFoundRequest(t *testing.T) {
	testEnv := testutils.SetupTestEnv()
	defer testEnv.TearDownTestEnv()
	testEnv.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		if r.Method != http.MethodGet {
			t.Errorf("got %s method, expected GET", r.Method)
		}
	})

	endpoint := testEnv.Server.URL + "/"
	client := &Client{
		HTTPClient: &http.Client{},
		Endpoint:   endpoint,
	}

	ctx := context.Background()
	response, err := client.doRequest(ctx, http.MethodGet, endpoint, nil)
	require.NoError(t, err)
	require.NotNil(t, response.Body)
	require.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestDoErrGenericRequest(t *testing.T) {
	testEnv := testutils.SetupTestEnv()
	defer testEnv.TearDownTestEnv()
	testEnv.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		require.Equal(t, http.MethodGet, r.Method)
	})

	endpoint := testEnv.Server.URL + "/"
	client := &Client{
		HTTPClient: &http.Client{},
		Endpoint:   endpoint,
	}

	ctx := context.Background()
	response, err := client.doRequest(ctx, http.MethodGet, endpoint, nil)
	require.NoError(t, err)
	require.NotNil(t, response.Body)
	require.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestDoErrNoContentRequest(t *testing.T) {
	testEnv := testutils.SetupTestEnv()
	defer testEnv.TearDownTestEnv()
	testEnv.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = fmt.Fprint(w, "") // write no content in the response body.

		require.Equal(t, http.MethodGet, r.Method)
	})

	endpoint := testEnv.Server.URL + "/"
	client := &Client{
		HTTPClient: &http.Client{},
		Endpoint:   endpoint,
	}

	ctx := context.Background()
	response, err := client.doRequest(ctx, http.MethodGet, endpoint, nil)
	require.NoError(t, err)
	require.NotNil(t, response.Body)
	require.Equal(t, http.StatusBadGateway, response.StatusCode)
}

func TestDoRequestInvalidResponseFromServer(t *testing.T) {
	testEnv := testutils.SetupTestEnv()
	defer testEnv.TearDownTestEnv()
	testEnv.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = fmt.Fprint(w, "<") // might be as a beginning of HTTP body

		if r.Method != http.MethodGet {
			t.Errorf("got %s method, want GET", r.Method)
		}
	})

	endpoint := testEnv.Server.URL + "/"
	client := &Client{
		HTTPClient: &http.Client{},
		Endpoint:   endpoint,
	}

	ctx := context.Background()
	response, err := client.doRequest(ctx, http.MethodGet, endpoint, nil)
	require.NoError(t, err)
	require.NotNil(t, response.Body)
	require.Equal(t, http.StatusServiceUnavailable, response.StatusCode)
}
