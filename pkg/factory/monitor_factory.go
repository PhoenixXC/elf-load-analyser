package factory

import (
    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
    "log"
)

var (
    factoryByName = make(map[string]*bcc.Monitor)
)

func Register(name string, monitor *bcc.Monitor)  {
    factoryByName[name] = monitor
}

func LoadMonitors(ctx bcc.Context)  {
    log.Printf("Load monitor...")
    for _, monitor := range factoryByName {
        _ = monitor.TouchOff(ctx.Pid)
        m, g := monitor.DoAction()
        if g {
            go monitor.Resolve(m)
        }
    }
}
