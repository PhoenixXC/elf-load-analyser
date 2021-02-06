package modules

import (
    "bytes"
    "unsafe"
)

func bytes2Str(arr []byte) string {
    l := bytes.IndexByte(arr[:], 0)
    arr = arr[:l]
    return *(*string)(unsafe.Pointer(&arr))
}