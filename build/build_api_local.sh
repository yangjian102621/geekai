#!/bin/bash

# build go api program
cd ../api

# get os type and arch

os=$(go env GOOS)
arch=$(go env GOARCH)

echo "CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -o bin/chatgpt-plus-api main.go"

CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -o bin/chatgpt-plus-api main.go