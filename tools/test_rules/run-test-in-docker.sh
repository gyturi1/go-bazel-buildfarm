#!/bin/bash
set -eu -o pipefail

DOCKER_INCREMENTAL_LODER=$1

#load docker image
sh -c $DOCKER_INCREMENTAL_LODER

#start the loaded container
docker run --rm -a stdout -a stdin -a stderr bazel/it-test:test_runner_image