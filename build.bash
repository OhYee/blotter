#!/bin/bash

DEPS_VERSION="0"
BASE_VERSION="0"
PROD_VERSION=$(git describe --abbrev=0 --tags 2>/dev/null || git rev-parse --short HEAD)

BUILD_DEPS=""
BUILD_BASE=""
BUILD_PROD="1"

for arg in "$@"; do 
    case "$arg" in
       "test") PROD_VERSION="test" ;;
       "deps") BUILD_DEPS="1" BUILD_BASE="" BUILD_PROD="" ;;
       "base") BUILD_DEPS="" BUILD_BASE="1" BUILD_PROD="" ;;
        esac
done

DEPS_IMAGE="ohyee/blotter:deps_${DEPS_VERSION}"
BASE_IMAGE="ohyee/blotter:base_${BASE_VERSION}"
PROD_IMAGE="ohyee/blotter:${PROD_VERSION}"

echo -e "DEPS VERSION:  $DEPS_IMAGE \t\t `if [[ -n \"$BUILD_DEPS\" ]]; then echo '√'; fi`"
echo -e "BASE VERSION:  $BASE_IMAGE \t\t `if [[ -n \"$BUILD_BASE\" ]]; then echo '√'; fi`"
echo -e "PROD VERSION:  $PROD_IMAGE \t\t `if [[ -n \"$BUILD_PROD\" ]]; then echo '√'; fi`"
echo ""

function build_docker() {
    echo "Building docker image $IMAGE from $DOCKER_FILE"
    docker build \
        --build-arg DEPS_IMAGE=$DEPS_IMAGE \
        --build-arg BASE_IMAGE=$BASE_IMAGE \
        --build-arg PROD_IMAGE=$PROD_IMAGE \
        -t $IMAGE \
        `if [[ -n $DOCKER_FILE ]]; then echo -n "-f $DOCKER_FILE"; fi` \
        .
}

if [[ -n $BUILD_PROD ]]; then
    IMAGE=$PROD_IMAGE
    build_docker
fi

if [[ -n $BUILD_BASE ]]; then
    IMAGE=$BASE_IMAGE
    DOCKER_FILE="Dockerfile.base"
    build_docker
fi

if [[ -n $BUILD_DEPS ]]; then
    IMAGE=$DEPS_IMAGE
    DOCKER_FILE="Dockerfile.deps"
    build_docker
fi
