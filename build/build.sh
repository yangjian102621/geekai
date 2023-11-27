#!/bin/bash

version=$1
# build go api
cd ../api
make clean linux

# build web app
cd ../web
npm run build

cd ../build

# remove docker image if exists
docker rmi -f registry.cn-shenzhen.aliyuncs.com/geekmaster/chatgpt-plus-api:$version
# build docker image for chatgpt-plus-go
docker build -t registry.cn-shenzhen.aliyuncs.com/geekmaster/chatgpt-plus-api:$version -f dockerfile-api-go ../

# build docker image for chatgpt-plus-vue
docker rmi -f registry.cn-shenzhen.aliyuncs.com/geekmaster/chatgpt-plus-web:$version
docker build --platform linux/amd64 -t registry.cn-shenzhen.aliyuncs.com/geekmaster/chatgpt-plus-web:$version -f dockerfile-vue ../

if [ "$2" = "push" ];then
  docker push registry.cn-shenzhen.aliyuncs.com/geekmaster/chatgpt-plus-api:$version
  docker push registry.cn-shenzhen.aliyuncs.com/geekmaster/chatgpt-plus-web:$version
fi
