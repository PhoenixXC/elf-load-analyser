package web

import (
    "fmt"
    "testing"
)

func TestGetAnyPort(t *testing.T) {
    fmt.Println(getAnyFreePort())
}
