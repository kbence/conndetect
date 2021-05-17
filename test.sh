#!/usr/bin/env bash

ROOT=$(cd "$(dirname ${BASH_SOURCE[0]:-$0})"; pwd)

cd "$ROOT"
go test $(find . -name '*_test.go' | sed 's/[^/]\+$//' | uniq)
