#!/bin/bash

TAG=$(git describe --abbrev=0 --tags 2>/dev/null || git rev-parse --short HEAD)

docker build -t ohyee/blotter:${TAG} .