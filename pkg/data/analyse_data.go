package data

import (
    "encoding/json"
    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/system"
    "log"
    "time"
)

type AnalyseData struct {
    success bool
    now     time.Time // time
    name    string    // event name
    data    string    // data by json
}

func NewAnalyseData(monitor *bcc.Monitor, data ...Data) *AnalyseData {
    jsonData := "{}"
    byteData, err := json.Marshal(data)
    if err != nil {
        log.Printf(system.Error("Convert %v to Json error, %v\n"), data, err)
    } else {
        jsonData = string(byteData[:])
    }
    return &AnalyseData{name: monitor.Name, success: err != nil, data: jsonData, now: time.Now()}
}
