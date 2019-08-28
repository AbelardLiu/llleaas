#!/bin/bash

IMAGE=llleaas
TAG=latest

docker inspect ${IMAGE}${TAG} &> /dev/null

if [ $? -ne 0 ]; then
    echo "not found image ${IMAGE}:${TAG}"
else
    docker image rm -f ${IMAGE}:${TAG}
fi