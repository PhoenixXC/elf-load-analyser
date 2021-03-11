package module

import (
    _ "embed" // for embed bcc source
    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/data/markdown"
    "github.com/phoenixxc/elf-load-analyser/pkg/modules"
    "github.com/phoenixxc/elf-load-analyser/pkg/modules/enhance"
)

//go:embed src/bprm_execve.cpp.k
var bprmExecveSrc string

type bprmExecveEvent struct {
    enhance.TimeEventResult
}

func (a bprmExecveEvent) Render() *data.AnalyseData {
    return data.NewAnalyseData("", markdown.NewTextContent("execve"))
}

type bprmExecve struct {
    modules.MonitorModule
}

func init() {
    m := modules.NewPerfResolveMonitorModule(&bprmExecve{})
    m.RegisterOnceTable("call_event", func(data []byte) (*data.AnalyseData, error) {
        return modules.Render(data, &bprmExecveEvent{}, true)
    })
    modules.ModuleDefaultInit(m)
}

func (a *bprmExecve) Monitor() string {
    return "bprm_execve"
}

func (a *bprmExecve) Source() string {
    return bprmExecveSrc
}

func (a *bprmExecve) Events() []*bcc.Event {
    ke := bcc.NewKprobeEvent("kprobe__bprm_execve", "bprm_execve", -1)
    return []*bcc.Event{ke}
}
