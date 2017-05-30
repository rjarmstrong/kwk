#!/usr/bin/env bash

set -ef -o pipefail

path=$(dirname $0)

go test ${path}/app/runtime -cover