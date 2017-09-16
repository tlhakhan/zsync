#!/bin/bash

[[ -d ./bin ]] && rm -rf ./bin && mkdir ./bin
export GOBIN=$PWD/bin

go get -insecure google.golang.org/grpc
go get -insecure golang.org/x/net/context
go get -insecure github.com/tlhakhan/golib/cmd
go get -insecure github.com/golang/protobuf/protoc-gen-go
go get -insecure github.com/golang/protobuf/proto

go install client.go
go install server.go
