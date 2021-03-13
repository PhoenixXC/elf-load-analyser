package xelf

import (
    "fmt"
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
    "testing"
)

func TestBuildDynamicInfo(t *testing.T) {
    f, err := getELFFile("/bin/ls")
    if err != nil {
        log.Error(err)
    }
    dynInfo, err := BuildDynamicInfo(f)
    if err != nil {
        log.Error(err)
    }
    fmt.Println(dynInfo)
}
