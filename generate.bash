#!/usr/bin/env bash

# 生成代码

_date=$(date '+%Y-%m-%d %H:%M:%S')
_branch=$(git rev-parse --abbrev-ref HEAD); _branch="${_branch}@"
_version=$(git describe --abbrev=0 --tags 2>/dev/null); if [[ -z $_version ]]; then _version="v0.0.0"; fi
_blotter_version="${_branch}${_version} (${_date})"

echo ${_blotter_version}

CGO_ENABLED=1 go build -ldflags "-X 'main._version=${_blotter_version}' -extldflags '-static -s -w -fpic'" 

unset _date _branch _version _blotter_version