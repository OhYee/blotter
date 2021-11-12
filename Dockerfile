# syntax=docker/dockerfile:experimental

ARG BASE_IMAGE=""
ARG DEPS_IMAGE=""

FROM ${DEPS_IMAGE} AS builder

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

FROM ${BASE_IMAGE} AS prod

WORKDIR /data/blotter

ENV mongoURI="mongodb:27017"

ENTRYPOINT [ "./blotter", "-address", "0.0.0.0:50000" ]

# gojieba 字典文件
COPY --from=builder /go/pkg/mod/github.com/ttys3/gojieba@v1.1.3/dict /go/pkg/mod/github.com/ttys3/gojieba@v1.1.3/dict

COPY --from=builder /data/blotter/blotter /data/blotter/blotter