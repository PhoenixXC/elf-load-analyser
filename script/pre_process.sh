#!/usr/bin/env bash
# used by golang run and debug

source ./merge.sh

cd "$(pwd)/../pkg/modules/module/src" || exit
merge_src
cd "$(pwd)/../pkg/web/templates" || exit
build_fronted