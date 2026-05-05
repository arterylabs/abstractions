#!/usr/bin/env sh
set -eu

mode=${1:?mode is required}
shift

if [ "$#" -eq 0 ]; then
  exit 0
fi

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repo_dir=$(dirname "$script_dir")
frontend_dir="$repo_dir/frontend"

frontend_files=""
for file in "$@"; do
  case "$file" in
    frontend/*)
      frontend_files="${frontend_files}
${file#frontend/}"
      ;;
  esac
done

if [ -z "$frontend_files" ]; then
  exit 0
fi

set --
old_ifs=$IFS
IFS='
'
for file in $frontend_files; do
  [ -n "$file" ] && set -- "$@" "$file"
done
IFS=$old_ifs

cd "$frontend_dir"

case "$mode" in
  eslint)
    pnpm exec eslint --fix -- "$@"
    ;;
  prettier)
    pnpm exec prettier --write -- "$@"
    ;;
  *)
    echo "Unsupported mode: $mode" >&2
    exit 1
    ;;
esac
