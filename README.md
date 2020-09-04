# bookish-spork

Simple Go application that provides HTTP API to in-memory cache

## Public API

- `/v1/set` - set value to cache

Example:
```bash
curl -i -X POST "127.0.0.1:63100/v1/set" -H "Content-Type: application/json" \
                                         -d '{"key": "some-key", "value": "some-value", "ttl": 10}'
HTTP/1.1 200 OK
Date: Fri, 04 Sep 2020 16:33:37 GMT
Content-Length: 0
```

If `ttl` is equal or less to 0 it means that the key will never get expired.

- `/v1/get/<key>` - get value from cache
```bash
curl -s -X GET "127.0.0.1:63100/v1/get/some-key" | json_pp
{
   "value" : "some-value"
}
```

- `/v1/keys` - get list of all keys in cache
```bash
curl -s -X GET "127.0.0.1:63100/v1/keys" | json_pp
{
   "keys" : [
      "some-key",
      "some-key-1"
   ]
}
```

- `/v1/remove/<key>` - remove key from cache
```bash
curl -i -X DELETE "127.0.0.1:63100/v1/remove/some-key"
HTTP/1.1 204 No Content
Date: Fri, 04 Sep 2020 16:38:33 GMT
```

- `/v1/rpush` - add value to a list or create a new one
```bash
curl -i  -X POST "127.0.0.1:63100/v1/rpush" -H "Content-Type: application/json" -d '{"key": "some-key", "value": "some-value", "ttl": 0}'
HTTP/1.1 200 OK
Date: Fri, 04 Sep 2020 16:42:13 GMT
Content-Length: 0

curl -i  -X POST "127.0.0.1:63100/v1/rpush" -H "Content-Type: application/json" -d '{"key": "some-key", "value": "some-value-1"}'
HTTP/1.1 200 OK
Date: Fri, 04 Sep 2020 16:43:01 GMT
Content-Length: 0

curl -s -X GET "127.0.0.1:63100/v1/get/some-key" | json_pp
{
   "value" : [
      "some-value",
      "some-value-1"
   ]
}
```

- `/v1/lindex/<key>/<list-index>` - get value by the index in list
Example:
```bash
curl -s -X GET "127.0.0.1:63100/v1/lindex/some-key/0" | json_pp
{
   "value" : "some-value"
}

curl -s -X GET "127.0.0.1:63100/v1/lindex/some-key/1" | json_pp
{
   "value" : "some-value-1"
}

curl -s -X GET "127.0.0.1:63100/v1/lindex/some-key/2" | json_pp
{
   "value" : null
}

curl -s -X GET "127.0.0.1:63100/v1/lindex/some-key/-2" | json_pp
{
   "error" : "index can't be negative"
}
```

- `/v1/hset` - add key-value pairs to hash map or create a new one
Example:
```bash
curl -i  -X POST "127.0.0.1:63100/v1/hset" -H "Content-Type: application/json"
                                           -d '{"key": "some-hm", "value": {"k0": "v0", "k1": "v1"}, "ttl": 0}'
HTTP/1.1 200 OK
Date: Fri, 04 Sep 2020 16:46:32 GMT
Content-Length: 0

curl -s -X GET "127.0.0.1:63100/v1/get/some-hm" | json_pp
{
   "value" : {
      "k1" : "v1",
      "k0" : "v0"
   }
}
```

- `/v1/hget/<key>/<hash-map-key>` - get a value by hash map key
Example:
```bash
curl -s -X GET "127.0.0.1:63100/v1/hget/some-hm/k0" | json_pp
{
   "value" : "v0"
}

curl -s -X GET "127.0.0.1:63100/v1/hget/some-hm/k1" | json_pp
{
   "value" : "v1"
}

curl -s -X GET "127.0.0.1:63100/v1/hget/some-hm/k3" | json_pp
{
   "value" : null
}
```


## Service API

Service API provides standard [pprof](https://golang.org/pkg/net/http/pprof/) endpoints.
By default, it listens at port 63101.

Example of fetching profile over HTTP:
```bash
go tool pprof http://127.0.0.1:63101/debug/pprof/profile
```

You could also visit `http://127.0.0.1:63101/debug/pprof/` in your browser and do some profiling.

## Build

Use the following command to build binary:

```bash
make build
```

Or you can build Docker image:

```bash
docker build -t bookish-spork .
```

## Running

Example of config file: [bookish-spork.yaml](bookish-spork.example.yaml)

Running locally:
```bash
./bookish-spork --config <path-to-yaml-config>
```

Running in Docker (requires `bookish-spork` image to be build, see the commands above):
```bash
docker run -p 63101:63101 -p 63100:63100 \
           -v (pwd)/bookish-spork.yaml:/etc/bookish-spork/bookish-spork.yaml \
           bookish-spork
```

Note, that running the command above you should have `bookish-spork.yaml` locally.

## Testing

Use the following command to run acceptance tests (you will need `docker-compose`):

```sh
make acc-tests
```

Use the following command to run only unit-tests and linters:

```sh
make tests
```

Use the following command to run only unit-tests:

```sh
make unittest
```

## Linters

Use the following command to run golangci-lint:

```sh
make golangci-lint
```

Use the following command to run golangci-lint with unit-tests:

```sh
make tests
```

## TODO

* Add benchmark tests
* Code refactoring of `http` and `qqcache` packages
* Add missing functions to work with `list` and `hash maps`
