package modules

import (
    "bytes"
    "unsafe"
)

type vmAreaStruct struct {
    vmStart uint64
    vmEnd   uint64
    vmFlags uint64
    vmPgoff uint64
}

type mmStruct struct {
    stackVm uint64
    totalVm uint64
}

func bytes2Str(arr []byte) string {
    l := bytes.IndexByte(arr[:], 0)
    arr = arr[:l]
    return *(*string)(unsafe.Pointer(&arr))
}
