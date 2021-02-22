package data

import (
    "bytes"
    "log"
    "time"
    "unsafe"
)

type Type int8

type Status int8

// Data Format Type
const (
    MarkdownType = Type(iota)
    GraphvizType
)

// Data Status
const (
    Success      = Status(iota) // success
    RuntimeError                // exec runtime error, such as kernel function return failed
)

var status2Desc = map[Status]string{
    Success:      "OK",
    RuntimeError: "Load process error at runtime",
}

type Data struct {
    Class Type
    Data  string
    Style string
}

func NewData(class Type, data string) *Data {
    return &Data{Class: class, Data: data}
}

func (d *Data) WithStyle(css string) {
    d.Style = css
}

type AnalyseData struct {
    status    Status
    desc      string
    timestamp time.Time // time
    name      string    // event name
    data      *Data     // data
    dataList  []*AnalyseData
    extra     map[string]string
}

func (a *AnalyseData) Desc() string {
    return a.desc
}

func (a *AnalyseData) Status() Status {
    return a.status
}

func (a *AnalyseData) Timestamp() time.Time {
    return a.timestamp
}

func (a *AnalyseData) Name() string {
    return a.name
}

func (a *AnalyseData) Data() *Data {
    return a.data
}

func (a *AnalyseData) DataList() []*AnalyseData {
    return a.dataList
}

func NewAnalyseData(name string, data *Data) *AnalyseData {
    return &AnalyseData{name: name, status: Success, data: data, desc: statusDesc(Success), timestamp: time.Now(),
        extra: map[string]string{}}
}

func NewListAnalyseData(name string, dataList []*AnalyseData) *AnalyseData {
    return &AnalyseData{name: name, status: Success, dataList: dataList, desc: statusDesc(Success),
        timestamp: time.Now(), extra: map[string]string{}}
}

func NewErrAnalyseData(name string, s Status, desc string) *AnalyseData {
    if s == Success {
        log.Fatalf("Error status cannnot be OK")
    }
    if len(desc) == 0 {
        desc = statusDesc(s)
    }
    return &AnalyseData{status: s, desc: desc, timestamp: time.Now(), name: name, extra: map[string]string{}}
}

func (a *AnalyseData) SetTime(t time.Time) *AnalyseData {
    a.timestamp = t
    return a
}

func (a *AnalyseData) SetDesc(desc string) *AnalyseData {
    a.desc = desc
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
        return "Unknown status"
    }
    return res
}

func Bytes2Str(arr []byte) string {
    return *(*string)(unsafe.Pointer(&arr))
}

func TrimBytes2Str(arr []byte) string {
    l := bytes.IndexByte(arr, 0)
    arr = arr[:l]
    return *(*string)(unsafe.Pointer(&arr))
}
