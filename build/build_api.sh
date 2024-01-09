#!/bin/bash

# build go api program
cd ../api
make clean macos

# remove docker image if exists
docker rmi -f chatgpt-plus-api:lastest
# build docker image for chatgpt-plus-go
docker build -t chatgpt-plus-api:lastest -f dockerfile-api-go ../