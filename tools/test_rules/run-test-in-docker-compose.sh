#!/bin/bash
set -eu -o pipefail

RUNNER_PATH=$1
DOCKER_INCREMENTAL_LODER=$2
COMPOSE_FILE1=$3
COMPOSE_FILE2=$4

#load test_runner_image
sh -c $DOCKER_INCREMENTAL_LODER

#compose up
FILE_OPTS="-f $COMPOSE_FILE1 -f $COMPOSE_FILE2"
ARGS="up -d"

docker compose $FILE_OPTS $ARGS
