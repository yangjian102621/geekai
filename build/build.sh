#!/bin/bash

# Docker 引擎：docker（默认，本机）或 wslc
# 用法：ENGINE=wslc ./build.sh <version> [arch] [push]
engine=${ENGINE:-docker}
if [ "$engine" = "wslc" ]; then
  DOCKER="wslc"
else
  DOCKER="docker"
fi
echo "使用 Docker 引擎: $DOCKER"

version=$1
arch=${2:-amd64}

# build go api program
cd ../api
rm -rf bin/geekai-api*
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/geekai-api main.go

# build website and administration console
cd ../web
pnpm run build

cd ../build
# remove docker image if exists
$DOCKER rmi -f 10.8.65.4/geekai/geekai-api:$version-$arch
# build docker image for Geek-AI API
$DOCKER build -t 10.8.65.4/geekai/geekai-api:$version-$arch -f dockerfile-api-go ../

# build docker image for Geek-AI-web
$DOCKER rmi -f 10.8.65.4/geekai/geekai-web:$version-$arch
$DOCKER build -t 10.8.65.4/geekai/geekai-web:$version-$arch -f dockerfile-web-$arch ../

if [ "$3" = "push" ];then
  $DOCKER push 10.8.65.4/geekai/geekai-api:$version-$arch
  $DOCKER push 10.8.65.4/geekai/geekai-web:$version-$arch
fi
