package render

import (
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
)

type Type string

const (
    ElfType     = Type("_Elf")
    EnvType     = Type("_Env")
    AnalyseType = Type("_Load")
)

type Render interface {
    Render() (*data.AnalyseData, error)
    Type() Type
}