# syntax=docker/dockerfile:experimental

FROM golang:1.17.3

WORKDIR /data/blotter

# deps cache
COPY ./go.mod ./go.sum /data/blotter/
RUN go mod download -x 
RUN go build all

# build code
# build with cache: https://github.com/golang/go/issues/27719
COPY ./ /data/blotter
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go generate