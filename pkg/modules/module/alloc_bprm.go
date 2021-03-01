package module

import (
    _ "embed"
    "fmt"
    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/data/markdown"
    "github.com/phoenixxc/elf-load-analyser/pkg/modules"
    "github.com/phoenixxc/elf-load-analyser/pkg/modules/enhance"
)

// about stack 512byte limit,
// see: https://stackoverflow.com/questions/53627094/ebpf-track-values-longer-than-stack-size
//go:embed src/alloc_bprm.cpp.k
var allowBprmSource string

type allocBprmEvent struct {
    enhance.TimeEventResult
    Filename    [256]byte
    Fdpath      [256]byte
    Interp      [256]byte
    CurTopOfMem uint64
    RlimCur     uint64
    RlimMax     uint64
}

func (a allocBprmEvent) Render() *data.AnalyseData {
    s := fmt.Sprintf("filename: %q, fdpath: %q, interp: %q, rlimit stack cur: 0x%X,"+
        " rlimit stack max: 0x%X, current of top mem: 0x%X",
        data.TrimBytes2Str(a.Filename[:]), data.TrimBytes2Str(a.Fdpath[:]), data.TrimBytes2Str(a.Interp[:]),
        a.RlimCur, a.RlimMax, a.CurTopOfMem)
    return data.NewAnalyseData("alloc_bprm", markdown.NewTextContent(s))
}

type allocBprm struct {
    modules.MonitorModule
}

func init() {
    m := modules.NewPerfResolveMonitorModule(&allocBprm{})
    m.RegisterOnceTable("events", func(data []byte) (*data.AnalyseData, error) {
        return modules.Render(data, &allocBprmEvent{}, true)
    })
    modules.ModuleDefaultInit(m)
}

func (a *allocBprm) Monitor() string {
    return "alloc_bprm"
}

func (a *allocBprm) Source() string {
    return allowBprmSource
}

func (a *allocBprm) Events() []*bcc.Event {
    ke := bcc.NewKretprobeEvent("kretprobe__alloc_bprm", "alloc_bprm", -1)
    return []*bcc.Event{ke}
}
