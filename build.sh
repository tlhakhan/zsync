#!/bin/bash

[[ -d ./bin ]] && rm -rf ./bin && mkdir ./bin
export GOBIN=$PWD/bin

go install client.go
go install server.go
