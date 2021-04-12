#!/bin/bash


SHELL_FOLDER=$(cd "$(dirname "$0")";pwd)
IMAGE="golang:1.16.3"


# 拉取镜像
if [[ $(docker images ${IMAGE} | wc -l) -eq "1" ]]; then
    echo "Pulling docker image ${IMAGE}..."
    docker pull ${IMAGE}
fi

# 更新代码
echo "Pulling latest code..."
git pull


# 在 docker 中挂载当前目录编译
echo "Building execute in docker..."

docker run --rm \
    -v ${SHELL_FOLDER}:/data/blotter \
    -v ${SHELL_FOLDER}/temp:/go/pkg \
    ${IMAGE} \
    bash -c "go env -w GOPROXY=https://goproxy.cn,direct && cd /data/blotter && echo 'Generating...' && go generate"

echo "Build finished"

docker build -t blotter .

echo "Docker image 'blotter' build finished."
echo 'Using `docker run --rm --name=backend blotter` to start server.'