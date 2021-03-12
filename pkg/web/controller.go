package web

import (
    "errors"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
    "github.com/phoenixxc/elf-load-analyser/pkg/render"
    "net"
    "net/http"
    "strconv"
    "syscall"
)

var (
    analyseDataCenter []*render.Data
)

// VisualAnalyseData 数据展示
func VisualAnalyseData(p *data.Pool, port uint) {
    renderedData := render.DoAnalyse(p)
    go startWebService(renderedData, port)
}

func startWebService(d []*render.Data, port uint) {
    analyseDataCenter = d
    http.Handle("/", FrontedService())
    http.HandleFunc("/api/report", AnalyseReportService)

    // 程序的主流程，若定义的端口被占用，使用随机端口
MAIN:
    addr, err := getAnyFreeAddr(port)
    if err != nil {
        log.Errorf("Cannot select port to start wev server: %v", err)
    }

    log.Infof(log.Emphasize("Try to start wev server on %s, "+
        "you can view analysis report through this link if the startup is successful"), "http://"+addr)
    err = http.ListenAndServe(addr, nil)
    if err != nil {
        errType := syscall.EADDRINUSE
        if port != 0 && errors.Is(err, errType) {
            log.Warnf("Failed to start web service on defined port: %d, %v", port, err)
            log.Warn("Retry use random port...")
            port = 0
            goto MAIN
        }
        log.Errorf("Failed to start web service, %v", err)
    }
}

func getAnyFreeAddr(port uint) (string, error) {
    if port != 0 {
        return "0.0.0.0:" + strconv.Itoa(int(port)), nil
    }

    var listener *net.TCPListener
    var addr *net.TCPAddr
    var err error

    if addr, err = net.ResolveTCPAddr("tcp", "0.0.0.0:0"); err == nil {
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
    return "0.0.0.0:" + strconv.Itoa(listener.Addr().(*net.TCPAddr).Port), nil
}
