#!/usr/bin/env bash

set -ef -o pipefail

path=$(dirname $0)

go test ${path}/runtime -cover
go test ${path}/updater -cover
go test ${path}/app/handlers -cover
go vet ${path}/app/...