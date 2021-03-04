package render

import "github.com/phoenixxc/elf-load-analyser/pkg/data"

type AnalyseRender struct {
    dataList []*data.AnalyseData
}

func NewAnalyseRender(dataList []*data.AnalyseData) *AnalyseRender {
    return &AnalyseRender{dataList: dataList}
}

func (a AnalyseRender) Render() (*data.AnalyseData, error) {
    t := a.Type()
    return data.NewListAnalyseData(t.Id, t.Name, a.dataList), nil
}

func (a AnalyseRender) Type() Type {
    return AnalyseType
}
