#!/usr/bin/env bash

set -euo pipefail

ROOT=$(cd "$(dirname ${BASH_SOURCE[0]:-$0})"; pwd)

cd "$ROOT"

go get -d .

for path in $(find . -name '*_test.go' | sed 's/[^/]\+$//' | sort | uniq); do
    go test -count 1 "$path"
done
