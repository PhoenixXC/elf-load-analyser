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

func Register(monitor *bcc.Monitor)  {
    factory = append(factory, monitor)
}

// LoadMonitors ctx The run context, ctr control the process when to stop
func LoadMonitors(ctx bcc.Context, ctr chan struct{})  {
    log.Printf("Load monitor...")
    p := data.NewPool()
    ch := p.Chan()
    p.Init()

    go func() {
        select {
        case <-ctr:
            time.Sleep(10 * time.Millisecond)
            p.Close()
        }
    }()

    // 关闭 chan 事件，中止所有的信息收集协程和监听数据池任务
    closeHandler := func() {
        close(ctr)
    }

    for _, monitor := range factory {
        _ = monitor.TouchOff(ctx.Pid)
        m, g := monitor.DoAction()
        if g {
            // TODO 抽象独立的程序中止 handler，添加 WithEnd 方法（设置优先级，因为有可能会因为内核不支持某个hook而取消），
            //  标记最后一个处理的收集程序，当执行最后一个收集器的结束方法后，关闭 ctr，隐藏 closeHandler
            go monitor.Resolve(m, ch, closeHandler)
        }
    }
}
