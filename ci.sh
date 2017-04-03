#!/usr/bin/env bash

docker build . --rm -t kwkcli --build-arg BUILD_NUMBER=${BUILDKITE_BUILD_NUMBER}
docker rmi -f kwkcli