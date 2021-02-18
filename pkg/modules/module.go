package modules

import (
    "bytes"
    "encoding/binary"
    "fmt"
    bpf "github.com/iovisor/gobpf/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/factory"
    "github.com/phoenixxc/elf-load-analyser/pkg/system"
    "log"
    "reflect"
    "strings"
)

var (
    defaultFlags []string
)

// EventResult 事件结果接口，实现中的类型大小在编译时必须是已知的
type EventResult interface {
    Render() *data.AnalyseData
}

// MonitorModule 模块抽象接口
type MonitorModule interface {
    // Monitor 返回模块的名称，以及是否作为结束标志
    Monitor() string
    // Source 返回注入的 bcc 源码
    Source() string
    // Events 返回要注册的事件
    Events() []*bcc.Event
    // Resolve 解析、发送处理结果
    Resolve(m *bpf.Module, ch chan<- *data.AnalyseData, ready chan<- struct{}, stop <-chan struct{})
    // Render 数据转换、渲染
    Render(data []byte, event EventResult) (*data.AnalyseData, error)
}

// ModuleInit 注册 Module
func ModuleInit(mm MonitorModule, end bool) {
    m := bcc.NewMonitor(mm.Monitor(), mm.Source(), defaultFlags, mm.Resolve)
    for _, event := range mm.Events() {
        m.AddEvent(event)
    }
    if end {
        m.SetEnd()
    }
    factory.Register(m)
}

// BaseMonitorModule 抽象实现，封装 MonitorModule
type BaseMonitorModule struct {
    MonitorModule
}

func (b *BaseMonitorModule) Render(data []byte, event EventResult) (*data.AnalyseData, error) {
    err := binary.Read(bytes.NewBuffer(data), bpf.GetHostByteOrder(), event)
    if err != nil {
        e := fmt.Errorf(system.Error("(%s, %s) Failed to decode received d: %v\n"),
            b.Monitor(), reflect.TypeOf(event).Name(), err)
        return nil, e
    }
    return event.Render(), nil
}

// tCtx table context
type tCtx struct {
    name string
    loop bool
    channel chan []byte
    handler func(data []byte) (*data.AnalyseData, error)
}

// PerfResolveMonitorModule BaseMonitorModule 的高级抽象，封装 table 和 resolve 的处理
type PerfResolveMonitorModule struct {
    *BaseMonitorModule
    tableIds  []string
    table2Ctx map[string]*tCtx
    stopHandler func(p *PerfResolveMonitorModule)
}

func NewPerfResolveMonitorModule(m MonitorModule) *PerfResolveMonitorModule {
    return &PerfResolveMonitorModule{
        BaseMonitorModule: &BaseMonitorModule{
            MonitorModule: m,
        },
        tableIds:  []string{},
        table2Ctx: map[string]*tCtx{},
        stopHandler: nil,
    }
}

func (p *PerfResolveMonitorModule) RegisterTable(name string, loop bool, handler func(data []byte) (*data.AnalyseData, error)) {
    name = strings.TrimSpace(name)
    if handler == nil || len(name) == 0 {
        return
    }
    p.tableIds = append(p.tableIds, name)
    p.table2Ctx[name] = &tCtx{
        name:    fmt.Sprintf("%s@%s", p.Monitor(), name),
        loop:    loop,
        channel: make(chan []byte),
        handler: handler,
    }
}

func (p *PerfResolveMonitorModule) RegisterStopHandle(handler func(p *PerfResolveMonitorModule)) {
    p.stopHandler = handler
}

func (p *PerfResolveMonitorModule) Resolve(m *bpf.Module, ch chan<- *data.AnalyseData,
    ready chan<- struct{}, stop <-chan struct{}) {
    if len(p.tableIds) == 0 {
        return
    }

    // init perf map
    perI := 0
    perfMaps := make([]*bpf.PerfMap, len(p.tableIds))
    for _, table := range p.tableIds {
        t := bpf.NewTable(m.TableId(table), m)
        perf, err := bpf.InitPerfMap(t, p.table2Ctx[table].channel, nil)
        if err != nil {
            log.Fatalf(system.Error("(%s, %s) Failed to init perf map: %v\n"), p.Monitor(), "events", err)
        }
        perfMaps[perI] = perf
        perI++
    }

    ok := make(chan struct{})

    go func() {
        defer func() { close(ok) }()

        chCnt := len(p.table2Ctx)
        cnt := chCnt + 1
        if p.stopHandler != nil {
            cnt++
        }
        remaining := cnt
        tableNames := make([]string, chCnt)
        cases := make([]reflect.SelectCase, cnt)

        i := 0
        for t, c := range p.table2Ctx {
            cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(c.channel)}
            tableNames[i] = t
            i++
        }
        cases[chCnt] = reflect.SelectCase{
            Dir:  reflect.SelectSend,
            Chan: reflect.ValueOf(ready),
            Send: reflect.ValueOf(struct{}{}),
        }
        if idx := chCnt + 1; idx < cnt {
            cases[idx] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(stop)}
        }

        for remaining > 0 {
            chosen, value, ok := reflect.Select(cases)
            if !ok {
                cases[chosen].Chan = reflect.ValueOf(nil)
                remaining -= 1
                continue
            }

            tName := tableNames[chosen]
            ctx := p.table2Ctx[tName]

            if chosen == chCnt {
                cases[chosen].Chan = reflect.ValueOf(nil)
            } else if chosen == chCnt+1 {
                p.stopHandler(p)
            } else {
                d := value.Bytes()
                log.Printf("Resolve %q...", ctx.name)
                analyseData, err := ctx.handler(d)
                if err != nil {
                    log.Printf("Event %q resolve error: %v", ctx.name, err)
                } else {
                    ch <- analyseData
                }

                if ctx.loop {
                    continue
                }
                cases[chosen].Chan = reflect.ValueOf(nil)
            }
            remaining -= 1
        }
    }()

    for _, perfMap := range perfMaps {
        perfMap.Start()
    }
    <-ok
    for _, perfMap := range perfMaps {
        perfMap.Stop()
    }
}
