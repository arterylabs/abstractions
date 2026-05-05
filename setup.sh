#!/usr/bin/env sh
set -eu

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
cd "$script_dir"

if ! command -v python3 >/dev/null 2>&1; then
  echo "python3 is required to install pre-commit." >&2
  exit 1
fi

venv_dir="$script_dir/.venv/pre-commit"

if [ ! -x "$venv_dir/bin/python" ]; then
  python3 -m venv "$venv_dir"
fi

"$venv_dir/bin/python" -m pip install --upgrade pip pre-commit

hooks_path=$(git config --local --get core.hooksPath || true)
case "$hooks_path" in
  "" )
    ;;
  .husky|.husky/_ )
    git config --local --unset core.hooksPath
    ;;
  * )
    echo "Cannot install pre-commit because core.hooksPath is set to '$hooks_path'." >&2
    echo "Unset or update that Git setting, then rerun ./setup.sh." >&2
    exit 1
    ;;
esac

"$venv_dir/bin/pre-commit" install --install-hooks --hook-type pre-commit
