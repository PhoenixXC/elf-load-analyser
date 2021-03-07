package web

import (
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
    "net"
    "net/http"
)

var (
    analyseDataCenter []*data.AnalyseData
)

func StartWebService(d []*data.AnalyseData) {
    analyseDataCenter = d

    addr, err := getAnyFreePort()
    if err != nil {
        log.Errorf("Cannot select port to start wev server: %v", err)
    }
    log.Infof(log.Emphasize("Start web service on %s, click to view Analyse Report"), "http://"+addr)

    http.Handle("/", FrontedService())
    http.HandleFunc("/api/analyse_report", AnalyseReportService)
    err = http.ListenAndServe(addr, nil)
    if err != nil {
        log.Errorf("Start web service failed, %v", err)
    }
}

func getAnyFreePort() (string, error) {
    addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
    if err != nil {
        return "", err
    }

    listener, err := net.ListenTCP("tcp", addr)
    if err != nil {
        return "", err
    }

    defer func() {
        e := listener.Close()
        if e != nil {
            log.Errorf("Release random port error, %v", e)
        }
    }()
    return listener.Addr().(*net.TCPAddr).String(), nil
}
