package web

import (
    "encoding/json"
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
    "net/http"
)

func AnalyseReportService(w http.ResponseWriter, _ *http.Request) {
    d, err := json.Marshal(analyseDataCenter)
    if err != nil {
        log.Errorf("Serialize report data failed, %v", err)
    }
    _, _ = w.Write(d)
}
