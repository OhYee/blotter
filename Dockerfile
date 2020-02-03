FROM ubuntu

ENV DOCKER_PROXY="http://host.docker.internal:1081"
ENV NODE_REGISTRY="https://registry.npm.taobao.org"
ENV GOPROXY="https://goproxy.cn,direct"

# set proxy and update
ENV HTTP_PROXY=$DOCKER_PROXY \
    HTTPS_PROXY=$DOCKER_PROXY \
    http_proxy=$DOCKER_PROXY \
    https_proxy=$DOCKER_PROXY 
ENV PATH=$PATH:/usr/local/bin

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates
RUN mv /etc/apt/sources.list /etc/apt/sources.list.back 
RUN printf 'deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic main restricted universe multiverse\ndeb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-updates main restricted universe multiverse\ndeb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-backports main restricted universe multiverse\ndeb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-security main restricted universe multiverse' >> /etc/apt/sources.list


RUN unset HTTP_PROXY
RUN unset HTTPS_PROXY
RUN unset http_proxy
RUN unset https_proxy

RUN apt-get update && apt-get install git screen golang nodejs npm yarn mongodb -y
RUN npm install -g yarn
RUN yarn config set registry $NODE_REGISTRY
RUN yarn config set proxy $DOCKER_PROXY
RUN service mongodb start
