package render

import "github.com/phoenixxc/elf-load-analyser/pkg/data"

type AnalyseRender struct {
    dataList []*data.AnalyseData
}

func NewAnalyseRender(dataList []*data.AnalyseData) *AnalyseRender {
    return &AnalyseRender{dataList: dataList}
}

func (a AnalyseRender) Render() (*data.AnalyseData, error) {
    return data.NewListAnalyseData(string(a.Type()), a.dataList), nil
}

func (a AnalyseRender) Type() Type {
    return AnalyseType
}



