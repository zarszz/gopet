#!/bin/bash

rm -rf pb/*.go
mkdir pb
protoc --experimental_allow_proto3_optional --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
  --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
  proto/*.proto
