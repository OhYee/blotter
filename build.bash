#!/bin/bash

TAG=$(git describe --abbrev=0 --tags 2>/dev/null || git rev-parse --short HEAD)

if [[ $1 == 'test' ]]; then
    TAG="test"
fi

IMAGE="ohyee/blotter:${TAG}"
echo $IMAGE
docker build -t ${IMAGE} .