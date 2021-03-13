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

func (e *EnvRender) Render() (*Data, error) {
    t := e.Type()
    envContent := markdown.NewTitleContents(markdown.H2, "内核").WithContents(env.GetKernelVersion()).
        Append(markdown.NewTitleContents(markdown.H2, "平台").WithContents(runtime.GOARCH)).
        Append(markdown.NewTitleContents(markdown.H2, "环境变量").Append(markdown.NewList(os.Environ()...)))
    return NewData(data.NewAnalyseData(t.Name, envContent).WithID(t.ID)), nil
}

func (e *EnvRender) Type() Type {
    return EnvType
}

func (e *EnvRender) Release() {}
