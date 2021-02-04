package bcc

import (
    bpf "github.com/iovisor/gobpf/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/system"
    "log"
    "strconv"
    "strings"
)

const (
    KprobesEvent    = 1 << iota // kprobes
    KretprobesEvent             // kretprobes
    SyscallEvent                // syscall
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

type RegisterHandler interface {
    // handle op, return is continue register and error
    Handle() (bool, error)
}

type Event struct {
    action
    RegisterHandler
    Class  int    // 事件类型
    Name   string // 事件名称
    FnName string // 函数名称
}

func NewEvent(class int, name string, fnName string) *Event {
    return &Event{Class: class, Name: name, FnName: fnName}
}

func (e *Event) WithHandle(handle RegisterHandler) {
    e.RegisterHandler = handle
}

type Monitor struct {
    isInit bool
    event2Action map[*Event]*action
    Name         string
    Source       string // 模块源
    CFlags       []string
    Resolve      func(m *bpf.Module, send chan<- data.AnalyseData, stop func())
}

func NewMonitor(name string, source string, cFlags []string,
    resolve func(m *bpf.Module, ch chan<- data.AnalyseData, stop func())) *Monitor {
    return &Monitor{Name: name, Source: source, CFlags: cFlags, Resolve: resolve}
}

// init 创建 bpf 模块
func (m *Monitor) init() *bpf.Module {
    m.isInit = true
    return bpf.NewModule(m.Source, m.CFlags)
}

// TouchOff pid 触发的默认实现
func (m *Monitor) TouchOff(pid int) error {
    m.Source = strings.Replace(m.Source, "_PID_", strconv.Itoa(pid), -1)
    return nil
}

// AddEvent 设置事件
func (m *Monitor) AddEvent(event *Event) *Monitor {
    // do pre Handle
    if event.RegisterHandler != nil {
        con, err := event.RegisterHandler.Handle()
        if err != nil {
            log.Printf(system.Error("Failed to register [%v] event, %v\n"), *event, err)
        }
        if !con {
            log.Printf(system.Warn("Ignore register [%v] event"), *event)
        }
    }

    if m.event2Action == nil {
        m.event2Action = make(map[*Event]*action)
    }
    m.event2Action[event] = &event.action
    return m
}

func (m *Monitor) DoAction() (*bpf.Module, bool) {
    module := m.init()
    goOn := false
    for event, action := range m.event2Action {
        fd, err := (*action).Load(module)
        if err != nil {
            log.Printf(system.Error("Failed to load event %v, %v\n"), *event, err)
            continue
        }
        if err = (*action).Attach(module, fd); err != nil {
            log.Printf(system.Error("Failed to attach event %v, %v\n"), *event, err)
            continue
        }
        goOn = true
    }
    return module, goOn
}
