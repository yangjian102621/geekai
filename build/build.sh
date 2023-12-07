#!/bin/bash

version=$1
arch=${2:-amd64}

# build go api program
cd ../api
make clean $arch

# build web app
cd ../web
npm run build

cd ../build

# remove docker image if exists
docker rmi -f registry.cn-shenzhen.aliyuncs.com/geekmaster/chatgpt-plus-api:$version-$arch
# build docker image for chatgpt-plus-go
docker build -t registry.cn-shenzhen.aliyuncs.com/geekmaster/chatgpt-plus-api:$version-$arch -f dockerfile-api-go ../

# build docker image for chatgpt-plus-vue
docker rmi -f registry.cn-shenzhen.aliyuncs.com/geekmaster/chatgpt-plus-web:$version-$arch
docker build --platform linux/amd64 -t registry.cn-shenzhen.aliyuncs.com/geekmaster/chatgpt-plus-web:$version-$arch -f dockerfile-vue ../

if [ "$3" = "push" ];then
  docker push registry.cn-shenzhen.aliyuncs.com/geekmaster/chatgpt-plus-api:$version-$arch
  docker push registry.cn-shenzhen.aliyuncs.com/geekmaster/chatgpt-plus-web:$version-$arch
fi
