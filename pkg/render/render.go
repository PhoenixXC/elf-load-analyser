package render

import (
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
)

type Type struct {
    Id   string
    Name string
}

var (
    ElfType     = Type{Id: "_ELF", Name: "文件格式"}
    EnvType     = Type{Id: "_ENV", Name: "系统环境"}
    AnalyseType = Type{Id: "_LOAD", Name: "加载过程"}
)

type Render interface {
    Render() (*data.AnalyseData, error)
    Type() Type
}

type Content struct {
    Filepath string
}
