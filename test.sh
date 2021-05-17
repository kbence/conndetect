#!/usr/bin/env bash

set -euo pipefail

ROOT=$(cd "$(dirname ${BASH_SOURCE[0]:-$0})"; pwd)

cd "$ROOT"

go get -d .

for path in $(find . -name '*_test.go' | sed 's/[^/]\+$//' | uniq); do
    go test "$path"
done
