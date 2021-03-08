package web

import (
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
    "net"
    "net/http"
    "strconv"
    "time"
)

var (
    analyseDataCenter []*data.AnalyseData
)

func StartWebService(d []*data.AnalyseData, port uint) {
    analyseDataCenter = d

    addr, err := getAnyFreeAddr(port)
    if err != nil {
        log.Errorf("Cannot select port to start wev server: %v", err)
    }
    go func() {
        // NOTE sleep to boot http serve boot finish
        time.Sleep(10 * time.Millisecond)
        log.Infof(log.Emphasize("Start web service on %s, click to view Analyse Report"), "http://"+addr)
    }()

    http.Handle("/", FrontedService())
    http.HandleFunc("/api/report", AnalyseReportService)
    err = http.ListenAndServe(addr, nil)
    if err != nil {
        log.Errorf("Start web service failed, %v", err)
    }
}

func getAnyFreeAddr(port uint) (string, error) {
    if port != 0 {
        return "0.0.0.0:" + strconv.Itoa(int(port)), nil
    }

    var listener *net.TCPListener
    var addr *net.TCPAddr
    var err error

    if addr, err = net.ResolveTCPAddr("tcp", "0.0.0.0:0"); err != nil {
        listener, err = net.ListenTCP("tcp", addr)
    }
    if err != nil {
        return "", err
    }

    defer func() {
        if listener == nil {
            return
        }
        e := listener.Close()
        if e != nil {
            log.Errorf("Release random port error, %v", e)
        }
    }()
    return listener.Addr().(*net.TCPAddr).String(), nil
}
