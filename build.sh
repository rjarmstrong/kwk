#!/usr/bin/env bash

set -ef -o pipefail

go test ./app

env GOOS=linux GOARCH=amd64 go build -x -o /builds/linux
env GOOS=darwin GOARCH=amd64 go build -x -o /builds/mac
env GOOS=windows GOARCH=amd64 go build -x -o /builds/windows