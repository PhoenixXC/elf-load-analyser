package modules

import (
    "fmt"
    bpf "github.com/iovisor/gobpf/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
    "reflect"
    "strings"
    "sync"
)

var (
    mutex              sync.Once
    registeredEnhancer = make(map[string]Enhancer)
)

func RegisteredEnhancer(name string, e Enhancer) {
    registeredEnhancer[name] = e
}

// TableCtx table context
type TableCtx struct {
    Name    string
    Monitor MonitorModule

    loop    bool
    channel chan []byte
    handler func(data []byte) (*data.AnalyseData, error)
    mark    map[string]struct{}
}

func (t *TableCtx) IsMark(mk string) bool {
    _, ok := t.mark[mk]
    return ok
}

// Enhancer enhance on PerfResolveMonitorModule.Resolve
type Enhancer interface {
    PreHandle(tCtx *TableCtx)
    AfterHandle(tCtx *TableCtx, aData *data.AnalyseData, err error) (*data.AnalyseData, error)
}

// PerfResolveMonitorModule BaseMonitorModule 的高级抽象，封装 table 和 resolve 的处理
type PerfResolveMonitorModule struct {
    MonitorModule
    tableIds    []string
    table2Ctx   map[string]*TableCtx
    stopHandler func(p *PerfResolveMonitorModule)
}

func NewPerfResolveMonitorModule(m MonitorModule) *PerfResolveMonitorModule {
    return &PerfResolveMonitorModule{
        MonitorModule: m,
        tableIds:      []string{},
        table2Ctx:     map[string]*TableCtx{},
        stopHandler:   nil,
    }
}

// RegisterOnceTable 注册 table，仅执行一次操作
func (p *PerfResolveMonitorModule) RegisterOnceTable(name string,
    handler func(data []byte) (*data.AnalyseData, error)) {
    p.RegisterTable(name, false, handler)
}

// RegisterTable 注册 table, 若 loop 为 true，返回对应的 chan，否则返回 nil
func (p *PerfResolveMonitorModule) RegisterTable(name string, loop bool,
    handler func(data []byte) (*data.AnalyseData, error)) chan<- []byte {
    name = strings.TrimSpace(name)
    if handler == nil || len(name) == 0 {
        return nil
    }
    tableChannel := make(chan []byte)
    p.tableIds = append(p.tableIds, name)
    p.table2Ctx[name] = &TableCtx{
        Name:    fmt.Sprintf("%s@%s", p.Monitor(), name),
        loop:    loop,
        channel: tableChannel,
        handler: handler,
        mark:    map[string]struct{}{},
        Monitor: p,
    }
    if !loop {
        return nil
    }
    return tableChannel
}

func (p *PerfResolveMonitorModule) RegisterStopHandle(handler func(p *PerfResolveMonitorModule)) {
    p.stopHandler = handler
}

func (p *PerfResolveMonitorModule) SetMark(name string, mk string) *PerfResolveMonitorModule {
    ctx, ok := p.table2Ctx[name]
    if !ok {
        return p
    }
    ctx.mark[mk] = struct{}{}
    return p
}

//nolint:funlen
func (p *PerfResolveMonitorModule) Resolve(m *bpf.Module, ch chan<- *data.AnalyseData, ready chan<- struct{}, stop <-chan struct{}) {
    if log.ConfigLevel() == log.DLevel {
        mutex.Do(func() {
            l := len(registeredEnhancer)
            if l == 0 {
                return
            }
            enhancers, idx := make([]string, l), 0
            for s := range registeredEnhancer {
                enhancers[idx] = s
                idx++
            }
            log.Debugf("Enhancer %v be registered", enhancers)
        })
    }

    if len(p.tableIds) == 0 {
        return
    }

    perfMaps := initPerMaps(m, p)
    ok := make(chan struct{})

    go func() {
        defer func() { close(ok) }()

        chCnt := len(p.table2Ctx)
        cnt := chCnt + 1
        if p.stopHandler != nil {
            cnt++
        }
        remaining := cnt
        cases, tableNames := buildSelectCase(cnt, p.table2Ctx, ready, stop)

        for remaining > 0 {
            chosen, value, ok := reflect.Select(cases)
            if !ok {
                cases[chosen].Chan = reflect.ValueOf(nil)
                remaining--
                continue
            }

            tName := tableNames[chosen]
            ctx := p.table2Ctx[tName]

            if chosen == chCnt+1 {
                p.stopHandler(p)
                return
            } else if chosen != chCnt {
                d := value.Bytes()
                dataProcessing(d, ctx, ch)
                if ctx.loop {
                    continue
                }
            }
            cases[chosen].Chan = reflect.ValueOf(nil)
            remaining--
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

func dataProcessing(d []byte, ctx *TableCtx, ch chan<- *data.AnalyseData) {
    for name, handler := range registeredEnhancer {
        log.Debugf("%s pre handle for %q", name, ctx.Name)
        handler.PreHandle(ctx)
    }

    log.Infof("Resolve %q...", ctx.Name)
    analyseData, err := ctx.handler(d)
    log.Debugf("Receive data from %q, %v", ctx.Name, analyseData)

    for name, handler := range registeredEnhancer {
        log.Debugf("%s after handle for %q", name, ctx.Name)
        analyseData, err = handler.AfterHandle(ctx, analyseData, err)
    }

    if err != nil {
        log.Warnf("Event %q resolve error: %v", ctx.Name, err)
    } else {
        if len(analyseData.Name) == 0 {
            analyseData.Name = ctx.Name
        }
        ch <- analyseData
    }
}

func buildSelectCase(cnt int, table2Ctx map[string]*TableCtx,
    ready chan<- struct{}, stop <-chan struct{}) ([]reflect.SelectCase, []string) {
    chCnt := len(table2Ctx)
    cases := make([]reflect.SelectCase, cnt)
    tableNames := make([]string, chCnt)

    i := 0
    for t, c := range table2Ctx {
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
    return cases, tableNames
}

func initPerMaps(m *bpf.Module, p *PerfResolveMonitorModule) []*bpf.PerfMap {
    perI := 0
    perfMaps := make([]*bpf.PerfMap, len(p.tableIds))
    for _, table := range p.tableIds {
        t := bpf.NewTable(m.TableId(table), m)
        perf, err := bpf.InitPerfMap(t, p.table2Ctx[table].channel, nil)
        if err != nil {
            log.Errorf("(%s, %s) Failed to init perf map: %v", p.Monitor(), "events", err)
        }
        perfMaps[perI] = perf
        perI++
    }
    return perfMaps
}
