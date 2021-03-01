package render

import (
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/data/markdown"
    "github.com/phoenixxc/elf-load-analyser/pkg/env"
    "os"
    "runtime"
    "strings"
)

type EnvRender struct{}

func NewEnvRender() *EnvRender {
    return &EnvRender{}
}

func (e *EnvRender) Render() (*data.AnalyseData, error) {
    thisTitle := string(e.Type())

    sysItem := markdown.NewContent().WithTitle(markdown.H3Level, "系统").WithContents(env.GetSysOS())
    archItem := markdown.NewContent().WithTitle(markdown.H3Level, "平台").WithContents(runtime.GOARCH)
    environItem := markdown.NewContent().WithTitle(markdown.H3Level, "环境变量").WithContents(strings.Join(os.Environ(), ";"))
    envItem := markdown.NewContent().WithTitle(markdown.H2Level, "环境").
        Append(sysItem).
        Append(archItem).
        Append(environItem)
    return data.NewAnalyseData(thisTitle, data.newData(data.MarkdownType, envItem.ToMarkdown())), nil
}

func (e *EnvRender) Type() Type {
    return EnvType
}
