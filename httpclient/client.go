package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	// defaultHTTPTimeout represents the default timeout (in seconds) for HTTP requests.
	defaultHTTPTimeout = 120

	// defaultDialTimeout represents the default timeout (in seconds) for HTTP connection establishments.
	defaultDialTimeout = 60

	// defaultKeepaliveTimeout represents the default keep-alive period for an active network connection.
	defaultKeepaliveTimeout = 60

	// defaultMaxIdleConns represents the maximum number of idle (keep-alive) connections.
	defaultMaxIdleConns = 100

	// defaultIdleConnTimeout represents the maximum amount of time an idle (keep-alive) connection will remain
	// idle before closing itself.
	defaultIdleConnTimeout = 100

	// defaultTLSHandshakeTimeout represents the default timeout (in seconds) for TLS handshake.
	defaultTLSHandshakeTimeout = 60

	// defaultExpectContinueTimeout represents the default amount of time to wait for a server's first
	// response headers.
	defaultExpectContinueTimeout = 1
)

const (
	getEndpoint    = "get"
	setEndpoint    = "set"
	keysEndpoint   = "keys"
	removeEndpoint = "remove"
	rpushEndpoint  = "rpush"
	lindexEndpoint = "lindex"
	hgetEndpoint   = "hget"
	hsetEndpoint   = "hset"
)

// Client stores details that are needed to work with bookish-spork.
type Client struct {
	// HTTPClient represents an initialized HTTP client that will be used to do requests.
	HTTPClient *http.Client

	// Endpoint represents an endpoint that will be used in all requests.
	Endpoint string
}

// NewClient initializes a new client for bookish-spork API.
func NewClient(endpoint string) *Client {
	return &Client{
		HTTPClient: newHTTPClient(),
		Endpoint:   endpoint,
	}
}

// NewClientCustomHTTP initializes a new client for bookish-spork API.
// using custom HTTP client.
// If customHTTPClient is nil - default HTTP client will be used.
func NewClientCustomHTTP(customHTTPClient *http.Client, endpoint string) *Client {
	if customHTTPClient == nil {
		customHTTPClient = newHTTPClient()
	}

	return &Client{
		HTTPClient: customHTTPClient,
		Endpoint:   endpoint,
	}
}

// newHTTPClient returns a reference to an initialized and configured HTTP client.
func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   defaultHTTPTimeout * time.Second,
		Transport: newHTTPTransport(),
	}
}

// newHTTPTransport returns a reference to an initialized and configured HTTP transport.
func newHTTPTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   defaultDialTimeout * time.Second,
			KeepAlive: defaultKeepaliveTimeout * time.Second,
		}).DialContext,
		MaxIdleConns:          defaultMaxIdleConns,
		IdleConnTimeout:       defaultIdleConnTimeout * time.Second,
		TLSHandshakeTimeout:   defaultTLSHandshakeTimeout * time.Second,
		ExpectContinueTimeout: defaultExpectContinueTimeout * time.Second,
	}
}

// doRequest performs the HTTP request with the current Client's HTTPClient.
// Authentication and optional headers will be added automatically.
func (client *Client) doRequest(ctx context.Context, method, path string, body io.Reader) (*ResponseResult, error) {
	// Prepare an HTTP request with the provided context.
	request, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	request = request.WithContext(ctx)

	// nolint
	// Send the HTTP request and populate the ResponseResult.
	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}

	responseResult := &ResponseResult{
		Response: response,
	}

	// Check status code and populate custom error body with extended error message if it's possible.
	if response.StatusCode >= http.StatusBadRequest {
		err = responseResult.extractErr()
	}
	if err != nil {
		return nil, err
	}

	return responseResult, nil
}

// ResponseResult represents a result of an HTTP request.
// It embeds standard http.Response and adds custom API error representations.
type ResponseResult struct {
	*http.Response

	*ErrNotFound

	*ErrGeneric

	// Err contains an error that can be provided to a caller.
	Err error
}

// ErrNotFound represents 404 status code error of an HTTP response.
type ErrNotFound struct {
	Error string `json:"error"`
}

// ErrGeneric represents a generic error of an HTTP response.
type ErrGeneric struct {
	Error string `json:"error"`
}

// ExtractResult allows to provide an object into which ResponseResult body will be extracted.
func (result *ResponseResult) extractResult(to interface{}) error {
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}
	defer result.Body.Close()

	return json.Unmarshal(body, to)
}

// extractErr populates an error message and error structure in the ResponseResult body.
func (result *ResponseResult) extractErr() error {
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}
	defer result.Body.Close()

	if len(body) == 0 {
		result.Err = fmt.Errorf("got the %d status code from the server", result.StatusCode)

		return nil
	}

	switch result.StatusCode {
	case http.StatusNotFound:
		err = json.Unmarshal(body, &result.ErrNotFound)
	default:
		err = json.Unmarshal(body, &result.ErrGeneric)
	}
	if err != nil {
		result.Err = fmt.Errorf("got invalid response from the server, status code %d",
			result.StatusCode)

		return nil
	}

	result.Err = fmt.Errorf("got the %d status code from the server: %s",
		result.StatusCode, string(body))

	return nil
}
