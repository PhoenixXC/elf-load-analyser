#!/bin/bash
# used by golang run and debug

source ./merge.sh

cd "$(pwd)/../pkg/modules/module/src" || exit
merge_src