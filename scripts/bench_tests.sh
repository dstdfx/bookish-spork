#!/usr/bin/env bash

echo "==> Running benchmark tests..."
go test -v -bench=. -run=^Benchmark ./...
if [[ $? -ne 0 ]]; then
    echo ""
    echo "Benchmark tests failed."
    exit 1
fi

exit 0