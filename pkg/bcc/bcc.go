package bcc

import (
    bpf "github.com/iovisor/gobpf/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
    "strconv"
    "strings"
)

const (
    KprobesEvent    = 1 << iota // kprobes
    KretprobesEvent             // kretprobes
)

type Context struct {
    Pid int
}

type action interface {
    // Attached symbol
    Attach(m *bpf.Module, fd int) error
    // Loader load symbol
    Load(m *bpf.Module) (int, error)
}

type Event struct {
    action
    Class  int    // 事件类型
    Name   string // 事件名称
    FnName string // 函数名称
}

func NewEvent(class int, name string, fnName string) *Event {
    return &Event{Class: class, Name: name, FnName: fnName}
}

type Monitor struct {
    isEnd        bool // end monitor level
    event2Action map[*Event]*action
    Name         string
    Source       string // 模块源
    CFlags       []string
    Resolve      func(m *bpf.Module, send chan<- *data.AnalyseData, ready chan<- struct{}, ok <-chan struct{})
}

func (m *Monitor) IsEnd() bool {
    return m.isEnd
}

func (m *Monitor) SetEnd() {
    m.isEnd = true
}

func NewMonitor(name string, source string, cFlags []string,
    resolve func(m *bpf.Module, ch chan<- *data.AnalyseData, ready chan<- struct{}, ok <-chan struct{})) *Monitor {
    return &Monitor{Name: name, Source: source, CFlags: cFlags, Resolve: resolve, isEnd: false}
}

// initialize 创建 bpf 模块
func (m *Monitor) initialize() *bpf.Module {
    return bpf.NewModule(m.Source, m.CFlags)
}

// PreProcessing 预处理
func (m *Monitor) PreProcessing(ctx Context) error {
    // PID replace
    m.Source = strings.ReplaceAll(m.Source, "_PID_", strconv.Itoa(ctx.Pid))
    return nil
}

// AddEvent 设置事件
func (m *Monitor) AddEvent(event *Event) *Monitor {
    if m.event2Action == nil {
        m.event2Action = make(map[*Event]*action)
    }
    m.event2Action[event] = &event.action
    return m
}

// DoAction 执行 attach operation 操作
func (m *Monitor) DoAction() (*bpf.Module, bool) {
    module := m.initialize()
    goOn := false
    for event, action := range m.event2Action {
        fd, err := (*action).Load(module)
        if err != nil {
            log.Warnf("Failed to load event %v, %v", *event, err)
        }
        if err == nil {
            if err = (*action).Attach(module, fd); err != nil {
                log.Warnf("Failed to attach event %v, %v", *event, err)
            }
        }
        if err != nil {
            if m.IsEnd() {
                log.Errorf("The necessary monitor %q start failed", m.Name)
            }
            continue
        }
        goOn = true
    }
    return module, goOn
}
