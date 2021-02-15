package data

import (
    "encoding/json"
    "github.com/phoenixxc/elf-load-analyser/pkg/system"
    "log"
    "time"
)

type Type int8

const (
    MarkdownType = Type(iota)
    GraphvizType
)

type Data struct {
    Class Type // Data type, MarkdownType or GraphvizType
    Data  string
    Style string // css style
}

func NewData(class Type, data string) *Data {
    return &Data{Class: class, Data: data}
}

func (d *Data) WithStyle(css string) {
    d.Style = css
}

type AnalyseData struct {
    success bool
    now     time.Time // time
    name    string    // event name
    data    string    // data
}

func (a *AnalyseData) Name() string {
    return a.name
}

func (a *AnalyseData) Data() string {
    return a.data
}

func NewAnalyseData(name string, data *Data) *AnalyseData {
    jsonData := "{}"
    byteData, err := json.Marshal(*data)
    if err != nil {
        log.Printf(system.Error("Convert %v to Json error, %v\n"), data, err)
    } else {
        jsonData = string(byteData[:])
    }
    return &AnalyseData{name: name, success: err != nil, data: jsonData, now: time.Now()}
}
