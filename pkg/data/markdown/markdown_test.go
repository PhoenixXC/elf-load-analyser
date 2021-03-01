package markdown

import (
    "fmt"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/env"
    "os"
    "runtime"
    "testing"
)

func TestItem(t *testing.T) {
    sysItem := NewContent().WithTitle(H3Level, "系统").WithContents(env.GetSysOS())
    archItem := NewContent().WithTitle(H3Level, "平台").WithContents(runtime.GOARCH)
    environItem := NewContent().WithTitle(H3Level, "环境变量").Append(NewList(os.Environ()...))
    envItem := NewContent().WithTitle(H2Level, "环境").Append(sysItem).Append(archItem).Append(environItem)
    fmt.Println(envItem.ToMarkdown())
}

func TestEmptyTitle(t *testing.T) {
    var content data.Builder = NewContent().WithContents("Test")
    fmt.Println(content.Data())
    fmt.Println(content.Class() == data.MarkdownType)
}

func TestTable(t *testing.T) {
    tb := NewTable("ID", "Name", "Sex", "Phone")
    tb.WithDesc("student record")
    tb.AddRow("1", "Alice")
    tb.AddRow("2", "Bob", "boy", "222222", "extra")
    tb.AddRow("3", "Tom", "girl", "3333333", "extra")
    tb.AddRow("4", "David", "boy", "444444")
    fmt.Println(tb.ToMarkdown())
}
