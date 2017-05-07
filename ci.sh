#!/usr/bin/env bash

docker build . --rm -t kwkcli --build-arg BUILD_NUMBER=${TRAVIS_JOB_ID}  AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
docker rmi -f kwkcli