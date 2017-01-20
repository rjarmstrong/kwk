#!/usr/bin/env bash

# Iterate directories and run protoc
for D in */; do /bin/bash -c "protoc -I $D --go_out=plugins=grpc:$D $D*.proto"; done

#https://github.com/google/protobuf/tree/master/js NOT DOING THIS - just importing proto file into node
#for D in */; do /bin/bash -c "protoc -I $D --js_out=library_style=commonjs,binary --grpc_out=$D #--plugin=protoc-gen-grpc=grpc_node_plugin $D*.proto"; done