package render

import (
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
)

type AnalyseRender struct {
    dataList []*data.AnalyseData
}

func NewAnalyseRender(dataList []*data.AnalyseData) *AnalyseRender {
    return &AnalyseRender{dataList: dataList}
}

func (a AnalyseRender) Render() (*Data, error) {
    t := a.Type()
    return NewData(data.NewListAnalyseData(t.ID, t.Name, a.dataList)), nil
}

func (a AnalyseRender) Type() Type {
    return AnalyseType
}

func (a AnalyseRender) Release() {}
