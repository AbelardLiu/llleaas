#!/bin/sh

CUR_PATH=`pwd`

COMPILE_SRC="${CUR_PATH}/compile/"
COMPILE_DST="/root/"

LLLEAAS_SOURCE_SRC="${CUR_PATH}/../"
LLLEAAS_SOURCE_DST="/go/src/lll.github.com/llleaas"

BUILD_SCRIPT="${COMPILE_DST}/build.sh"
BUILD_IMAGE="golang:alpine"

# build binary in container
docker run -v ${COMPILE_SRC}:${COMPILE_DST} -v ${LLLEAAS_SOURCE_SRC}:${LLLEAAS_SOURCE_DST} -i ${BUILD_IMAGE} ${BUILD_SCRIPT}

# mv binary into dst directory
APP_NAME="llleaas"
APP_SRC="${LLLEAAS_SOURCE_SRC}/cmd/hermes/${APP_NAME}"
APP_DST_NAME="hermes"
APP_DST="${CUR_PATH}/docker/_output/bin"

# mk dst dir
mkdir -p ${APP_DST}
mv ${APP_SRC} ${APP_DST}/${APP_DST_NAME}
