package factory

import (
    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
    "sync"
    "time"
)

var (
    mutex   sync.Mutex
    factory []*bcc.Monitor
)

func Register(monitor *bcc.Monitor) {
    factory = append(factory, monitor)
}

// LoadMonitors ctx The run context, ctr control the process when to stop
func LoadMonitors(ctx bcc.Context) (p *data.Pool, ok <-chan struct{}) {
    ready, o := make(chan struct{}), make(chan struct{})
    p, ok = data.NewPool(), o
    ch := p.Chan()
    p.Init()

    go func() {
        <-ok
        time.Sleep(10 * time.Millisecond)
        p.Close()
    }()

    log.Info("Start load monitor....")
    cnt := 0
    endFlag := false
    for _, monitor := range factory {
        flag := endFlag
        monitor := monitor
        err := monitor.PreProcessing(ctx)
        if err != nil {
            logHandle := log.Warnf
            if monitor.IsEnd() && !flag {
                logHandle = log.Errorf
            }
            logHandle("Monitor %q pre processing error: %v", monitor.Name, err)
            continue
        }
        if m, g := monitor.DoAction(); g {
            cnt++
            go func(o chan struct{}) {
                log.Infof("Monitor %q start...", monitor.Name)
                monitor.Resolve(m, ch, ready, ok)
                if monitor.IsEnd() && !flag {
                    close(o)
                    endFlag = true
                }
                log.Debugf("Close monitor: %q", monitor.Name)
                mutex.Lock()
                defer mutex.Unlock()
                m.Close()
            }(o)
        }
    }

    for ; cnt > 0; cnt-- {
        <-ready
    }
    log.Info("Load monitors ok")

    return
}
