package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Get returns value by key in cache.
func (client *Client) Get(ctx context.Context, key string) (interface{}, *ResponseResult, error) {
	url := strings.Join([]string{client.Endpoint, getEndpoint, key}, "/")
	responseResult, err := client.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	// Extract response body
	var v struct {
		Value interface{} `json:"value"`
	}

	err = responseResult.extractResult(&v)
	if err != nil {
		return nil, responseResult, err
	}

	return v.Value, responseResult, nil
}

// SetBody represents set request body.
type SetBody struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	TTL   int         `json:"ttl"`
}

// Get returns value by key in cache.
func (client *Client) Set(ctx context.Context, body SetBody) (*ResponseResult, error) {
	url := strings.Join([]string{client.Endpoint, setEndpoint}, "/")
	v, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	responseResult, err := client.doRequest(ctx, http.MethodPost, url, bytes.NewReader(v))
	if err != nil {
		return nil, err
	}
	if responseResult.Err != nil {
		return responseResult, responseResult.Err
	}

	return responseResult, nil
}

// Keys returns slice of all keys in cache.
func (client *Client) Keys(ctx context.Context) ([]string, *ResponseResult, error) {
	url := strings.Join([]string{client.Endpoint, keysEndpoint}, "/")
	responseResult, err := client.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	// Extract response body
	var v struct {
		Value []string `json:"keys"`
	}

	err = responseResult.extractResult(&v)
	if err != nil {
		return nil, responseResult, err
	}

	return v.Value, responseResult, nil
}

// Get returns value by key in cache.
func (client *Client) Remove(ctx context.Context, key string) (*ResponseResult, error) {
	url := strings.Join([]string{client.Endpoint, removeEndpoint, key}, "/")
	responseResult, err := client.doRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	if responseResult.Err != nil {
		return responseResult, responseResult.Err
	}

	return responseResult, nil
}

// RPushBody represents rpush request body.
type RPushBody struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	TTL   int         `json:"ttl"`
}

// RPush adds value to a list or create a new one.
func (client *Client) RPush(ctx context.Context, body RPushBody) (*ResponseResult, error) {
	url := strings.Join([]string{client.Endpoint, rpushEndpoint}, "/")
	v, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	responseResult, err := client.doRequest(ctx, http.MethodPost, url, bytes.NewReader(v))
	if err != nil {
		return nil, err
	}
	if responseResult.Err != nil {
		return responseResult, responseResult.Err
	}

	return responseResult, nil
}

// LIndex returns value by index in list.
func (client *Client) LIndex(ctx context.Context, key string, index int) (interface{}, *ResponseResult, error) {
	url := strings.Join([]string{client.Endpoint, lindexEndpoint, key, strconv.Itoa(index)}, "/")
	responseResult, err := client.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	// Extract response body
	var v struct {
		Value interface{} `json:"value"`
	}

	err = responseResult.extractResult(&v)
	if err != nil {
		return nil, responseResult, err
	}

	return v.Value, responseResult, nil
}

// HSetBody represents hset request body.
type HSetBody struct {
	Key   string                 `json:"key"`
	Value map[string]interface{} `json:"value"`
	TTL   int                    `json:"ttl"`
}

// HSet adds key-value pairs to hash map or create a new one.
func (client *Client) HSet(ctx context.Context, body HSetBody) (*ResponseResult, error) {
	url := strings.Join([]string{client.Endpoint, hsetEndpoint}, "/")
	v, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	responseResult, err := client.doRequest(ctx, http.MethodPost, url, bytes.NewReader(v))
	if err != nil {
		return nil, err
	}
	if responseResult.Err != nil {
		return responseResult, responseResult.Err
	}

	return responseResult, nil
}

// HGet returns a value by hash map key.
func (client *Client) HGet(ctx context.Context, key, hkey string) (interface{}, *ResponseResult, error) {
	url := strings.Join([]string{client.Endpoint, hgetEndpoint, key, hkey}, "/")
	responseResult, err := client.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	// Extract response body
	var v struct {
		Value interface{} `json:"value"`
	}

	err = responseResult.extractResult(&v)
	if err != nil {
		return nil, responseResult, err
	}

	return v.Value, responseResult, nil
}
