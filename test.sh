#!/bin/bash

export COMPOSE_PROJECT_NAME="$(basename $(pwd))-${RANDOM}"

remove_resources() {
    docker-compose down
}

trap remove_resources EXIT

docker-compose build
docker-compose run tests "${@}"
