FROM golang:1.16.3 AS builder

COPY ./ /data/blotter

WORKDIR /data/blotter

RUN go get
RUN go generate

# FROM golang:1.16.3 AS prod
FROM ubuntu AS prod

RUN apt update && \
    apt install -y --no-install-recommends \
    ca-certificates \
    graphviz \
    python3 \
    plantuml \
    python3-pip && \
    python3 -m pip install matplotlib && \
    rm -rf /var/lib/apt/lists/*

ENV mongoURI="mongodb:27017"

COPY --from=builder /data/blotter/blotter /data/blotter/blotter

# # gojieba 字典文件
COPY --from=builder /go/pkg/mod/github.com/ttys3/gojieba@v1.1.3/dict/hmm_model.utf8 /go/pkg/mod/github.com/ttys3/gojieba@v1.1.3/dict/hmm_model.utf8
COPY --from=builder /go/pkg/mod/github.com/ttys3/gojieba@v1.1.3/dict/idf.utf8 /go/pkg/mod/github.com/ttys3/gojieba@v1.1.3/dict/idf.utf8
COPY --from=builder /go/pkg/mod/github.com/ttys3/gojieba@v1.1.3/dict/jieba.dict.utf8 /go/pkg/mod/github.com/ttys3/gojieba@v1.1.3/dict/jieba.dict.utf8
COPY --from=builder /go/pkg/mod/github.com/ttys3/gojieba@v1.1.3/dict/stop_words.utf8 /go/pkg/mod/github.com/ttys3/gojieba@v1.1.3/dict/stop_words.utf8
COPY --from=builder /go/pkg/mod/github.com/ttys3/gojieba@v1.1.3/dict/user.dict.utf8 /go/pkg/mod/github.com/ttys3/gojieba@v1.1.3/dict/user.dict.utf8



WORKDIR /data/blotter

ENTRYPOINT [ "./blotter", "-address", "0.0.0.0:50000" ]