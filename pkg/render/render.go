package render

import (
    "fmt"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "time"
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
    Render() (*Data, error)
    Type() Type
}

type Content struct {
    Filepath string
}

/* Reformat render data structure */

type JsonTime time.Time

func (t JsonTime) MarshalJSON() ([]byte, error) {
    stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("15:04:05.000000"))
    return []byte(stamp), nil
}

type Data struct {
    ID       string   `json:"id"`
    Name     string   `json:"name"`
    Status   int      `json:"status"`
    Desc     string   `json:"desc"`
    GTime    JsonTime `json:"time"`
    Data     string   `json:"data"`
    GType    int      `json:"type"`
    DataList []*Data  `json:"dataList"`
}

func NewData(d *data.AnalyseData) *Data {
    renderData := &Data{
        ID:     d.ID,
        Name:   d.Name,
        Status: int(d.Status),
        Desc:   d.Desc,
        GTime:  JsonTime(d.Timestamp),
    }
    if d.Data != nil {
        renderData.Data = d.Data.Data
        renderData.GType = int(d.Data.Class)
    }
    if d.DataList != nil {
        renderData.DataList = []*Data{}
        for _, analyseData := range d.DataList {
            renderData.DataList = append(renderData.DataList, NewData(analyseData))
        }
    }
    return renderData
}
