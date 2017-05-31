#!/usr/bin/env bash

docker run \
    -e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
    -e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
    -v $PWD:/go/src/github.com/kwk-super-snippets/cli \
    --rm \
    rjarmstrong/goaws ./build.sh $(git log --pretty=format:'%h' -n 1)