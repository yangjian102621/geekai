#!/bin/bash

options=("build api" "build web")

select opt in "${options[@]}"
do
    case $opt in
        "build api")
            echo "you chose build api"
            ./build_api_local.sh
            ;;
        "build web")
            echo "you chose build web"
            ./build_web.sh
            ;;
        *) echo "invalid option $REPLY";;
    esac
    break
done