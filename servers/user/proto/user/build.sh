#!/bin/bash

#required the following
#protoc
#protoc-gen-go
#protoc-gen-micro

protoc --proto_path=$GOPATH/src:. --micro_out=. --go_out=. *.proto
