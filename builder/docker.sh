#!/bin/bash

IMAGE=llleaas
TAG=latest

cp -r ../config docker/

cd docker/
docker build -t ${IMAGE}:${TAG} -f Dockerfile .

rm -rf config