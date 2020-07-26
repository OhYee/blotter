module github.com/OhYee/blotter

go 1.13

require (
	github.com/OhYee/auth_github v1.0.3
	github.com/OhYee/auth_qq v1.0.1
	github.com/OhYee/goldmark-dot v1.0.2
	github.com/OhYee/goldmark-fenced_codeblock_extension v1.0.0
	github.com/OhYee/goldmark-image v1.0.0
	github.com/OhYee/goldmark-plantuml v1.0.2
	github.com/OhYee/goutils v1.0.1
	github.com/OhYee/rainbow v1.0.3
	github.com/alecthomas/chroma v0.7.3 // indirect
	github.com/gorilla/schema v1.1.0
	github.com/gorilla/websocket v1.4.2
	github.com/graemephi/goldmark-qjs-katex v0.2.0
	github.com/kisielk/errcheck v1.2.0 // indirect
	github.com/klauspost/compress v1.10.10 // indirect
	github.com/mitchellh/mapstructure v1.3.3
	github.com/qiniu/api.v7 v7.2.5+incompatible
	github.com/qiniu/api.v7/v7 v7.5.0
	github.com/yanyiwu/gojieba v1.1.2
	github.com/yuin/goldmark v1.2.0
	github.com/yuin/goldmark-highlighting v0.0.0-20200307114337-60d527fdb691
	go.mongodb.org/mongo-driver v1.3.5
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899 // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208 // indirect

)

replace github.com/yanyiwu/gojieba v1.1.2 => github.com/ttys3/gojieba v1.1.3
