#!/bin/sh

CUR_PATH=`pwd`

COMPILE_SRC="${CUR_PATH}/compile/"
COMPILE_DST="/root/"

LLLEAAS_SOURCE_SRC="${CUR_PATH}/../"
LLLEAAS_SOURCE_DST="/go/src/lll.github.com/llleaas"

BUILD_SCRIPT="${COMPILE_DST}/build.sh"
BUILD_IMAGE="golang:alpine"

# mv binary into dst directory
APP_NAME="llleaas"
APP_SRC="${LLLEAAS_SOURCE_SRC}/cmd/${APP_NAME}"
APP_DST="${CUR_PATH}/docker/_output/bin"

if [ -d ${APP_DST} ]; then
    rm -rf ${APP_DST}
fi