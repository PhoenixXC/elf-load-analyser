package data

import (
    "fmt"
    "github.com/phoenixxc/elf-load-analyser/pkg/env"
    "os"
    "runtime"
    "strings"
    "testing"
)

func TestItem(t *testing.T) {
    sysItem := NewItem(H3Level, "系统", env.GetSysOS())
    archItem := NewItem(H3Level, "平台", runtime.GOARCH)
    environItem := NewItem(H3Level, "环境变量", strings.Join(os.Environ(), ";"))
    envItem := NewItem(H2Level, "环境", sysItem.String(), archItem.String(), environItem.String())
    fmt.Println(envItem.String())
}

func TestTable(t *testing.T) {
    tb := NewTable("ID", "Name", "Sex", "Phone")
    tb.WithDesc("student record")
    tb.AddRow("1", "Alice")
    tb.AddRow("2", "Bob", "boy", "222222", "extra")
    tb.AddRow("3", "Tom", "girl", "3333333", "extra")
    tb.AddRow("4", "David", "boy", "444444")
    fmt.Println(tb.String())
}
