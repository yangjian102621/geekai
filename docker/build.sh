#!/bin/bash

version=$1
# build go api
cd ../api
make clean linux

# build web app
cd ../web
npm run build

cd ../docker

# remove docker image if exists
docker rmi -f registry.cn-hangzhou.aliyuncs.com/geekmaster/chatgpt-plus-go:$version
# build docker image for chatgpt-plus-go
docker build -t registry.cn-hangzhou.aliyuncs.com/geekmaster/chatgpt-plus-go:$version -f dockerfile-api-go ../

# build docker image for chatgpt-plus-vue
docker rmi -f registry.cn-hangzhou.aliyuncs.com/geekmaster/chatgpt-plus-vue:$version
docker build --platform linux/amd64 -t registry.cn-hangzhou.aliyuncs.com/geekmaster/chatgpt-plus-vue:$version -f dockerfile-vue ../

