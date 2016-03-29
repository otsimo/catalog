#!/bin/bash

export IMPORT_PATH=$GOPATH/src:./protos
export GENERATOR="gogofaster_out"
export OUTPUT_DIR="."
export PROTO_FILES="./protos/*.proto"

protoc --proto_path=$IMPORT_PATH --${GENERATOR}=plugins=grpc:${OUTPUT_DIR} $PROTO_FILES