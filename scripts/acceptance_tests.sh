#!/usr/bin/env bash

teardown() {
    docker-compose down --remove-orphans
    if [[ $? -ne 0 ]]; then
        echo ""
        echo "Acceptance tests cleanup failed. Check for dangling containers and networks."
        exit 1
    fi
    cd "${wayback}"
}

echo "==> Running acceptance tests with basic unit tests..."
wayback="$(pwd)"
tests_dir="${wayback}/test/acceptance"
cd ${tests_dir}
docker-compose up --abort-on-container-exit bookish-spork
if [[ $? -ne 0 ]]; then
    echo ""
    echo "Acceptance tests suite failed."
    teardown
    exit 1
fi

echo "Acceptance tests suite succeeded."
teardown
exit 0
