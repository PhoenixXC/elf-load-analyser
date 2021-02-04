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
func LoadMonitors(ctx bcc.Context, ctr chan struct{}) *data.Pool {
    log.Printf("Load monitor...")
    p := data.NewPool()
    ch := p.Chan()
    p.Init()

    go func() {
        if <-ctr; true {
            time.Sleep(10 * time.Millisecond)
            p.Close()
        }
    }()

    for _, monitor := range factory {
        _ = monitor.TouchOff(ctx.Pid)
        m, g := monitor.DoAction()
        if g {
            go func() {
                monitor.Resolve(m, ch)
                if monitor.IsEnd() {
                    close(ctr)
                }
            }()
        }
    }
    return p
}
