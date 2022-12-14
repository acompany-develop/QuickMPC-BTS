#!/bin/bash

build() {
    docker-compose build small_bts
}

setup() {
    docker-compose -f docker-compose.yml down -v
}

run() {
    docker-compose run small_bts /bin/sh -c 'cd /QuickMPC-BTS && go mod vendor && go test -count 1 -cover -v $(go list github.com/acompany-develop/QuickMPC-BTS/src/BeaverTripleService/Server/EngineToBts/...)'
}

teardown() {
    docker-compose -f docker-compose.yml down -v
}
