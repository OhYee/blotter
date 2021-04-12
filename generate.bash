#!/usr/bin/env bash

# 生成代码

_date=$(date '+%Y-%m-%d %H:%M:%S')
_branch=$(git rev-parse --abbrev-ref HEAD); _branch="${_branch}@"
_version=$(git describe --abbrev=0 --tags 2>/dev/null); if [[ -z $_version ]]; then _version="v0.0.0"; fi
_ldflags="${_branch}${_version} (${_date})"

go build -ldflags "-X 'main._version=${_ldflags}'"

unset _date _branch _version _ldflags