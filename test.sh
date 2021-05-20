#!/usr/bin/env bash

set -euo pipefail

ROOT=$(cd "$(dirname ${BASH_SOURCE[0]:-$0})"; pwd)

cd "$ROOT"

go get -d .
go test ./...
