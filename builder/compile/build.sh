#!/bin/sh

CPATH=`pwd`
ROOTPATH="${CPATH}/../../"

APP_NAME="llleaas"

GOPATH="/go/"
WORKPATH="${GOPATH}/src/lll.github.com/llleaas/cmd/hermes"

# add git
# apk update
# apk add git

# go get deps
# go get github.com/gorilla/mux

# build binary
cd ${WORKPATH}
CGO_ENABLED=0 GOOS=linux go build -o ${APP_NAME} -a -ldflags '-extldflags "-static"' .