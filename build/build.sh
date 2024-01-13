#!/bin/bash

echo 'Please enter your choice: '

options=("build local bin" "build docker" "build docker and push all")

select opt in "${options[@]}"
do
    case $opt in
        "build local bin")
            echo "you chose build local bin"
            bash ./build_local_bin.sh
            ;;
        "build docker")
            echo "you chose build docker"
            bash ./docker_build.sh
            ;;
        "build docker and push all")
            echo "you chose build docker and push all"
            bash ./docker_build.sh push
            ;;
        *) echo "invalid option $REPLY";;
    esac
    break
done
