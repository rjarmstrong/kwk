#!/usr/bin/env bash

# Iterate directories and run protoc
#for F in */; do /bin/bash -c "protoc -I $D --gofast_out=plugins=grpc:$D $D*.proto"; done

protoc --gofast_out=plugins=grpc:. *.proto

#protoc-gen-combo --gogo_out=plugins=grpc:. --proto_path=.:"${GOPATH}/src/:${GOPATH}/src/github.com/gogo/protobuf/protobuf/" users.proto
#protoc-gen-gogofaster --gogo_out=plugins=grpc:. --proto_path=.:"${GOPATH}/src/:${GOPATH}/src/github.com/gogo/protobuf/grpc.proto"