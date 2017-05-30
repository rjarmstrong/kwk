#!/bin/sh

set -ef -o pipefail

path=$(dirname $0)

go test ${path}/app/runtime -cover