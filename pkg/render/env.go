package render

import (
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
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

    sysItem := data.NewItem(data.H3Level, "系统", env.GetSysOS())
    archItem := data.NewItem(data.H3Level, "平台", runtime.GOARCH)
    environItem := data.NewItem(data.H3Level, "环境变量", strings.Join(os.Environ(), ";"))
    envItem := data.NewItem(data.H2Level, "环境", sysItem.String(), archItem.String(), environItem.String())
    return data.NewAnalyseData(thisTitle, data.NewData(data.MarkdownType, envItem.String())), nil
}

func (e *EnvRender) Type() Type {
    return EnvType
}
