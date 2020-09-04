# bookish-spork

Simple Go application that provides HTTP API to in-memory cache

## Public API

wip

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

## Usage

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

