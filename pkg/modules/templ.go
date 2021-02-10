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
    "reflect"
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
