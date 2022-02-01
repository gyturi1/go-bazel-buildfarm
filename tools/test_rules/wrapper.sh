#!/bin/bash

set -u -o pipefail

#start app
/cmd &> /tmp/app.log &

#wait for startup
until cat /tmp/app.log | grep -q -m 1 "listening on port"; do echo "waiting for server to spin up" sleep 0.3; done

#run tests
/it-test_test