# bookish-spork http client

This directory contains simple HTTP API client in Go for bookish-spork in-memory cache server

_NOTE! It's located in the same repository with the server just for the sake of simplicity._

Example usage:
```go
package main

import (
	"context"
	"fmt"

	bs "github.com/dstdfx/bookish-spork/httpclient"
)

func main() {
	// Init client
	cli := bs.NewClient("http://0.0.0.0:63100/v1")

	// Init common context
	ctx := context.Background()

	// Set key-value to cache
	_, err := cli.Set(ctx, bs.SetBody{
		Key:   "some-key",
		Value: "some-value",
		TTL:   10, // TTL in seconds
	})
	if err != nil {
		panic(err)
	}

	// Get value from cache
	v, _, err := cli.Get(ctx, "some-key")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", v)
}
```