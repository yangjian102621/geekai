#!/bin/bash

version=$1
arch=${2:-amd64}

# build go api program
cd ../api
make clean $arch

# build website and administration console
cd ../web
npm run build

cd ../build

# remove docker image if exists
docker rmi -f 10.8.65.4/geekai/geekai-api:$version-$arch
# build docker image for Geek-AI API
docker build -t 10.8.65.4/geekai/geekai-api:$version-$arch -f dockerfile-api-go ../

# build docker image for Geek-AI-web
docker rmi -f 10.8.65.4/geekai/geekai-web:$version-$arch
docker build -t 10.8.65.4/geekai/geekai-web:$version-$arch -f dockerfile-web-$arch ../

if [ "$3" = "push" ];then
  docker push 10.8.65.4/geekai/geekai-api:$version-$arch
  docker push 10.8.65.4/geekai/geekai-web:$version-$arch
fi
