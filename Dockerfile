FROM ubuntu:latest

ENV mongoURI="mongodb:27017"

COPY ./blotter /blotter

# gojieba 字典文件
COPY ./temp/mod/github.com/ttys3/gojieba@v1.1.3/dict/hmm_model.utf8 /go/pkg/mod/github.com/ttys3/gojieba@v1.1.3/dict/hmm_model.utf8
COPY ./temp/mod/github.com/ttys3/gojieba@v1.1.3/dict/idf.utf8 /go/pkg/mod/github.com/ttys3/gojieba@v1.1.3/dict/idf.utf8
COPY ./temp/mod/github.com/ttys3/gojieba@v1.1.3/dict/jieba.dict.utf8 /go/pkg/mod/github.com/ttys3/gojieba@v1.1.3/dict/jieba.dict.utf8
COPY ./temp/mod/github.com/ttys3/gojieba@v1.1.3/dict/stop_words.utf8 /go/pkg/mod/github.com/ttys3/gojieba@v1.1.3/dict/stop_words.utf8
COPY ./temp/mod/github.com/ttys3/gojieba@v1.1.3/dict/user.dict.utf8 /go/pkg/mod/github.com/ttys3/gojieba@v1.1.3/dict/user.dict.utf8


ENTRYPOINT [ "/blotter", "-address=0.0.0.0:50000"]