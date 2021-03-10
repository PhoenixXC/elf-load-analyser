#!/usr/bin/env bash
# used by golang run and debug

source ./merge.sh
readonly THIS_DIR="$(pwd)"

cd "${THIS_DIR}/../pkg/modules/module/src" || exit
merge_src
cd "${THIS_DIR}/../pkg/web/template" || exit
build_fronted 'use cache'