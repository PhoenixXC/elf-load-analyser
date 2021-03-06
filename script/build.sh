#!/usr/bin/env bash
set -u
source ./util.sh
source ./merge.sh

readonly TARGET="ela"
readonly WORK_DIR=$(pwd)
readonly GO_WORK="${WORK_DIR}/../"
readonly SRC_DIR="${WORK_DIR}/../pkg/modules/module/src"
readonly FRONTED_DIR="${WORK_DIR}/../pkg/web/template"

readonly ARG_ERR=100
readonly VAL_ERR=101

rm_extra_symbol=0
compress_level=0

ld_flags=
gc_flags=
target_bin="${TARGET}"

function help() {
    echo "Usage of build.sh"
    echo "  -r bool"
    echo "     remove debug and symbol data"
    echo "  -l int"
    echo "     compressed level, 1-9; if define, compressed target file used by upx"
    echo "  -d bool"
    echo "     enable debug"
}

function get_arg_val() {
    if [ $# -lt 2 ] || [ -z "$2" ]; then
        echo "miss arg[$1] value"
        help
        exit "${ARG_ERR}"
    fi
    echo "$2"
}

while (($# > 0)); do
    case $1 in
    -r)
        rm_extra_symbol=1
        shift 1
        ;;
    -d)
        gc_flags="all=-N -l"
        shift 1
        ;;
    -l)
        level=$(get_arg_val "$@")
        shift 2
        if [ "${level}" -gt 9 ] || [ "${level}" -lt 1 ]; then
            echo "argue value invalid"
            help
            exit "${VAL_ERR}"
        fi
        compress_level="${level}"
        ;;
    *)
        echo "Invalid arg: $1"
        help
        exit "${ARG_ERR}"
    esac
done

if [ "${gc_flags}" ]; then
    echo "Debug mode, disable other feature"
    rm_extra_symbol=0
    compress_level=0
fi

echo "Bcc source dir: ${SRC_DIR}"

cd "${SRC_DIR}" || exit
echo "Merge files"
merge_src
echo "Build fronted"
cd "${FRONTED_DIR}" || exit
build_fronted

cd "${GO_WORK}" || exit
echo "Build binary..."
mkdir -p "target"
if [ "${rm_extra_symbol}" -ne 0 ]; then
    ld_flags="-s -w"
fi
go build -gcflags="${gc_flags}" -ldflags "${ld_flags}" -o target/"${TARGET}" github.com/phoenixxc/elf-load-analyser/cmd

# compressed
if [ "${compress_level}" -ne 0 ]; then
    cd target || exit
    target_bin="${TARGET}-compressed"
    echo "Clean old binary"
    rm_file "${target_bin}"

    echo "Compressed..."
    upx "-${compress_level}" -o "${target_bin}" "${TARGET}"

    rm_file "${TARGET}"
    cd "${GO_WORK}" || exit
fi

echo "Build ok, now you can use '$(pwd)/target/${target_bin}' run program"

cd "${SRC_DIR}" || exit
echo "Clean tmp files"
rm_file "*.${OUT_SUFFIX}"
