package module

import (
    _ "embed"
    "fmt"
    bpf "github.com/iovisor/gobpf/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/modules"
    "github.com/phoenixxc/elf-load-analyser/pkg/modules/enhance"
)

type sysExecveEvent struct {
    enhance.TimeEventResult
}

func (e sysExecveEvent) Render() *data.AnalyseData {
    return data.NewAnalyseData("syscall:execve", data.NewData(data.MarkdownType, "start call"))
}

type sysExecveRetEvent struct {
    enhance.TimeEventResult
    Ret int8
}

func (s sysExecveRetEvent) Render() *data.AnalyseData {
    item := "syscall:execve"
    if s.Ret != 0 {
        return data.NewErrAnalyseData(item, data.RuntimeError, fmt.Sprintf("execve failed, return %d", s.Ret))
    }
    return data.NewAnalyseData(item, data.NewData(data.MarkdownType, "execve success"))
}

//go:embed src/execve.cpp.k
var execveSource string

type sysExecve struct {
    modules.MonitorModule
}

func init() {
    entry := "call_events"
    m := modules.NewPerfResolveMonitorModule(&sysExecve{})
    m.RegisterOnceTable(entry, func(d []byte) (*data.AnalyseData, error) {
        return modules.Render(d, &sysExecveEvent{}, true)
    })
    m.RegisterOnceTable("ret_events", func(d []byte) (*data.AnalyseData, error) {
        return modules.Render(d, &sysExecveRetEvent{}, true)
    })
    m.SetMark(entry, enhance.StartMark)
    modules.ModuleInit(m, true)
}

func (e *sysExecve) Monitor() string {
    return "syscall:execve"
}

func (e *sysExecve) Source() string {
    return execveSource
}

func (e *sysExecve) Events() []*bcc.Event {
    fnName := bpf.GetSyscallFnName("execve")
    k := bcc.NewKprobeEvent("syscall__execve", fnName, -1)
    ke := bcc.NewKretprobeEvent("do_ret_sys_execve", fnName, -1)
    return []*bcc.Event{ke, k}
}
