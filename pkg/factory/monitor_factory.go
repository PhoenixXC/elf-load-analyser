package factory

import (
    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "log"
    "time"
)

var (
    factory []*bcc.Monitor
)

func Register(monitor *bcc.Monitor) {
    factory = append(factory, monitor)
}

// LoadMonitors ctx The run context, ctr control the process when to stop
func LoadMonitors(ctx bcc.Context, ok chan struct{}) *data.Pool {
    ready := make(chan struct{})
    p := data.NewPool()
    ch := p.Chan()
    p.Init()

    go func() {
        <-ok
        time.Sleep(10 * time.Millisecond)
        p.Close()
    }()

    log.Println("Start load monitor....")
    cnt := 0
    for _, monitor := range factory {
        _ = monitor.TouchOff(ctx.Pid)
        m, g := monitor.DoAction()
        if g {
            cnt++
            go func() {
                log.Printf("Monitor %s start to resolve...\n", monitor.Name)
                monitor.Resolve(m, ch, ready, ok)
                if monitor.IsEnd() {
                    close(ok)
                }
                defer m.Close()
            }()
        }
    }

    for ; cnt > 0; cnt-- {
        <-ready
    }
    log.Println("Start load monitor ok")

    return p
}
