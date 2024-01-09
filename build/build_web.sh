#!/bin/bash

# build web app
cd ../web
npm run build

cd ../build

# build docker image for chatgpt-plus-vue
docker rmi -f chatgpt-plus-web:lastest
docker build -t chatgpt-plus-web:latest -f dockerfile-vue ../