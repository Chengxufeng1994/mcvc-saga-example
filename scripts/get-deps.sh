#!/bin/bash

go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.1
go install github.com/google/wire/cmd/wire@v0.6.0
go install github.com/golang/protobuf/protoc-gen-go@v1.5.4
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0