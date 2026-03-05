#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset

[[ -n "${DEBUG:-}" ]] && set -x

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__file="${__dir}/$(basename "${BASH_SOURCE[0]}")"

cd $__dir/../static/terraform

terraform providers lock \
    -platform=linux_arm64 \
    -platform=linux_amd64 \
    -platform=darwin_arm64 \
    -platform=darwin_amd64
