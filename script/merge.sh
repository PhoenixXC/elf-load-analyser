#!/bin/bash

readonly OUT_SUFFIX="k"

function merge_src() {
    ls -- *."${OUT_SUFFIX}" >/dev/null 2>&1 && rm -- *."${OUT_SUFFIX}"
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