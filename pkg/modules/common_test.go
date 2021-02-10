package modules

import (
    "bytes"
    "testing"
)

func Benchmark_bytes2Str_unsafe(b *testing.B) {
    data := []byte("51cee58d2009745b72acb79005a881e4")
    a := [256]byte{}
    for i, d := range data {
        a[i] = d
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = bytes2Str(a[:])
    }
}

func Benchmark_bytes2Str_string(b *testing.B) {
    data := []byte("51cee58d2009745b72acb79005a881e4")
    a := [256]byte{}
    for i, d := range data {
        a[i] = d
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = string(a[:bytes.IndexByte(a[:], 0)])
    }
}
