package render

import (
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/data/markdown"
    "github.com/phoenixxc/elf-load-analyser/pkg/env"
    "os"
    "runtime"
)

type EnvRender struct{}

func NewEnvRender() *EnvRender {
    return &EnvRender{}
}

func (e *EnvRender) Render() (*data.AnalyseData, error) {
    t := e.Type()

    envContent := markdown.NewTitleContents(markdown.H3, "系统").WithContents(env.GetSysOS()).
        Append(markdown.NewTitleContents(markdown.H3, "平台").WithContents(runtime.GOARCH)).
        Append(markdown.NewTitleContents(markdown.H3, "环境变量").Append(markdown.NewList(os.Environ()...)))
    return data.NewAnalyseData(t.Name, envContent).WithId(t.Id), nil
}

func (e *EnvRender) Type() Type {
    return EnvType
}
