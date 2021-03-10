#!/usr/bin/env bash

source ./util.sh

readonly OUT_SUFFIX="k"
readonly FRONTED_CHECKSUM="md5sums"

function merge_src() {
    rm_file "*.${OUT_SUFFIX}"
    for line in *.cpp
    do
        grep -E '^#include[[:blank:]]+"[[:print:]]+"$' "${line}" | awk '{print $2}' | tr -d '"' \
        | awk '{
            if (NF > 0) {
                if ($1!="_dev.h") print $1
            }}' | xargs cat > "${line}"."${OUT_SUFFIX}"
        echo "" >> "${line}"."${OUT_SUFFIX}"
        sed '/^[[:blank:]]*#include[[:blank:]]\+"[[:print:]]\+"$/d' "${line}" | cat >> "${line}"."${OUT_SUFFIX}"
    done
}

function get_md5sums() {
    find ./ -type f \
            -not -path "./node_modules/*"\
            -not -path ".*/.umi/*"\
            -not -path ".*/.idea/*" \
            -not -path "./${FRONTED_CHECKSUM}"\
            -exec md5sum {} \; | sort -k 2 | md5sum | awk '{print $1}'
}

function build_fronted() {
    if (( $# > 0 )) && [ -f "${FRONTED_CHECKSUM}" ] && [ -d "dist" ]; then
        old_checksum=$(cat "${FRONTED_CHECKSUM}")
        new_checksum=$(get_md5sums)
        echo "old_checksum: ${old_checksum}"
        echo "new_checksum: ${new_checksum}"
        if [ "${old_checksum}" == "${new_checksum}" ]; then
            echo "file not change"
            return
        fi
    fi

    echo "build fronted..."
    export MOCK=none
    yarn
    yarn build
    get_md5sums > "${FRONTED_CHECKSUM}"
}
