package data

import (
    "bytes"
    "unsafe"
)

func Bytes2Str(arr []byte) string {
    return *(*string)(unsafe.Pointer(&arr))
}

func TrimBytes2Str(arr []byte) string {
    l := bytes.IndexByte(arr, 0)
    arr = arr[:l]
    return *(*string)(unsafe.Pointer(&arr))
}
