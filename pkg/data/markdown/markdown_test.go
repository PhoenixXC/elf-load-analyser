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
    sysItem := NewTitleContents(H3, "系统").WithContents(env.GetSysOS())
    archItem := NewTitleContents(H3, "平台").WithContents(runtime.GOARCH)
    environItem := NewTitleContents(H3, "环境变量").Append(NewList(os.Environ()...))
    envItem := NewTitleContents(H2, "环境").Append(sysItem).Append(archItem).Append(environItem)
    fmt.Println(envItem.ToMarkdown())
}

func TestEmptyTitle(t *testing.T) {
    var content data.Builder = NewTextContent("Test")
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
