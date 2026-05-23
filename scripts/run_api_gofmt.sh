#!/usr/bin/env sh
set -eu

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
api_dir=$(dirname "$script_dir")/api

if [ "$#" -eq 0 ]; then
  gofmt -w "$api_dir"
  exit 0
fi

gofmt -w "$@"
