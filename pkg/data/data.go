package data

import (
    "log"
    "time"
)

type Type int8

type Status int8

// Data Format Type
const (
    MarkdownType = Type(iota + 1)
)

const (
    Success      = Status(iota) // success
    RuntimeError                // exec runtime error, such as kernel function return failed
)

var status2Desc = map[Status]string{
    Success:      "OK",
    RuntimeError: "happened error at runtime",
}

type Builder interface {
    Class() Type
    Data() string
}

type Data struct {
    Class Type   `json:"class"`
    Data  string `json:"data"`
}

func newData(b Builder) *Data {
    return &Data{
        Class: b.Class(),
        Data:  b.Data(),
    }
}

type AnalyseData struct {
    Id        string         `json:"id"`
    Name      string         `json:"name"`
    Status    Status         `json:"status"`
    Desc      string         `json:"desc"`
    Timestamp time.Time      `json:"timestamp"`
    Data      *Data          `json:"render_data"`
    DataList  []*AnalyseData `json:"render_data_list"`
    Style     string         `json:"style"`
    extra     map[string]string
}

func (a *AnalyseData) DataStr() string {
    if a.Data == nil {
        return ""
    }
    return a.Data.Data
}

func NewAnalyseData(name string, builder Builder) *AnalyseData {
    return &AnalyseData{Name: name, Status: Success, Data: newData(builder), Desc: statusDesc(Success),
        Timestamp: time.Now(), extra: map[string]string{}}
}

func NewListAnalyseData(id string, name string, dataList []*AnalyseData) *AnalyseData {
    return &AnalyseData{
        Id:        id,
        Name:      name,
        Status:    Success,
        DataList:  dataList,
        Desc:      statusDesc(Success),
        Timestamp: time.Now(),
        extra:     map[string]string{},
    }
}

func NewErrAnalyseData(name string, s Status, desc string) *AnalyseData {
    if s == Success {
        log.Fatalf("Error Status cannnot be OK")
    }
    if len(desc) == 0 {
        desc = statusDesc(s)
    }
    return &AnalyseData{Status: s, Desc: desc, Timestamp: time.Now(), Name: name, extra: map[string]string{}}
}

func (a *AnalyseData) WithId(id string) *AnalyseData {
    a.Id = id
    return a
}

func (a *AnalyseData) PutExtra(k string, v string) {
    a.extra[k] = v
}

func (a *AnalyseData) Extra(k string) (string, bool) {
    v, ok := a.extra[k]
    return v, ok
}

func statusDesc(s Status) string {
    res, ok := status2Desc[s]
    if !ok {
        return "Unknown Status"
    }
    return res
}
